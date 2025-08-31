package integtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/murilo-bracero/sequence-technical-test/internal/db"
	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/murilo-bracero/sequence-technical-test/internal/handlers"
	"github.com/murilo-bracero/sequence-technical-test/internal/repository"
	"github.com/murilo-bracero/sequence-technical-test/internal/server"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/config"
	"github.com/murilo-bracero/sequence-technical-test/internal/services"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type EnvironmentCommands struct {
	pgContainer *postgres.PostgresContainer
	db          db.DB
}

func New() *EnvironmentCommands {
	return &EnvironmentCommands{}
}

func (e *EnvironmentCommands) Start(ctx context.Context) error {
	err := e.startDB(ctx)
	if err != nil {
		return err
	}
	err = e.startApp(ctx)
	if err != nil {
		return err
	}

	// wait for app to be ready
	time.Sleep(1 * time.Second)
	return nil
}

func (e *EnvironmentCommands) CreateSequence(ctx context.Context, body dto.CreateSequenceRequest) (*dto.SequenceResponse, error) {
	url := "http://localhost:8000/sequences"

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	req.Close = true

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var responseBody dto.SequenceResponse
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (e *EnvironmentCommands) GetSequenceById(ctx context.Context, id string) (*dto.SequenceResponse, error) {
	url := fmt.Sprintf("http://localhost:8000/sequences/%s", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Close = true

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var responseBody dto.SequenceResponse
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (e *EnvironmentCommands) ClearDatabase(ctx context.Context) error {
	if e.db == nil {
		return fmt.Errorf("database not initialized")
	}

	tx, err := e.db.Tx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM sequences")
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (e *EnvironmentCommands) Destroy(ctx context.Context) error {
	if e.pgContainer != nil {
		return e.pgContainer.Terminate(ctx)
	}
	return nil
}

func (e *EnvironmentCommands) startDB(ctx context.Context) error {
	image := "postgres:17.6-alpine3.22"

	c, err := postgres.Run(ctx, image,
		postgres.WithUsername("integtestpostgres"),
		postgres.WithPassword("integtestpostgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		))
	if err != nil {
		return err
	}

	postgresConnStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return err
	}

	err = e.migrate("postgres", postgresConnStr)
	if err != nil {
		return err
	}

	postgresConnStr, err = e.changeDBName(postgresConnStr, "sequencemailbox")
	if err != nil {
		return err
	}

	err = e.migrate("sequencemailbox", postgresConnStr)
	if err != nil {
		return err
	}

	e.pgContainer = c
	return nil
}

func (e *EnvironmentCommands) startApp(ctx context.Context) error {
	host, err := e.pgContainer.Host(ctx)
	if err != nil {
		return err
	}

	port, err := e.pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return err
	}

	cfg := &config.Config{
		PostgresHost:     host,
		PostgresPort:     port.Int(),
		PostgresUser:     "integtestpostgres",
		PostgresPassword: "integtestpostgres",
		PostgresDatabase: "sequencemailbox",
		MaxDbConnections: 10,
		MinDbConnections: 1,
		MaxConnIdleTime:  30,
	}

	db, err := db.New(context.Background(), cfg)
	if err != nil {
		return err
	}

	e.db = db

	sequenceRepository := repository.NewSequenceRepository(db)

	sequenceService := services.NewSequenceService(sequenceRepository)

	sequenceHandler := handlers.NewSequenceHandler(cfg, sequenceService)

	stepRepository := repository.NewStepRepository(db)

	stepService := services.NewStepService(sequenceRepository, stepRepository)

	stepHandler := handlers.NewStepHandler(stepService)

	go server.Start(db, sequenceHandler, stepHandler)

	return nil
}

func (e *EnvironmentCommands) migrate(dbName string, conn string) error {
	m, err := migrate.New("file://../db/migrations/"+dbName, conn)
	if err != nil {
		return err
	}

	return m.Up()
}

func (e *EnvironmentCommands) changeDBName(connString, newDBName string) (string, error) {
	u, err := url.Parse(connString)
	if err != nil {
		return "", fmt.Errorf("failed to parse connection string: %w", err)
	}

	u.Path = path.Join("/", newDBName)
	return u.String(), nil
}
