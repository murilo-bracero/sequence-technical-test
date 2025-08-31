package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	dao "github.com/murilo-bracero/sequence-technical-test/internal/db/gen"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/config"
)

type DB interface {
	Queries() *dao.Queries
	Close()
	Tx(context.Context) (pgx.Tx, error)
	Ping(context.Context) error
}

type db struct {
	pool *pgxpool.Pool
}

func New(context context.Context, cfg *config.Config) (DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.MaxDbConnections)
	poolConfig.MinConns = int32(cfg.MinDbConnections)
	poolConfig.HealthCheckPeriod = 30 * time.Second
	poolConfig.MaxConnIdleTime = time.Duration(cfg.MaxConnIdleTime) * time.Second

	pool, err := pgxpool.NewWithConfig(context, poolConfig)
	if err != nil {
		return nil, err
	}

	return &db{pool: pool}, nil
}

func (d *db) Queries() *dao.Queries {
	return dao.New(d.pool)
}

func (d *db) Close() {
	d.pool.Close()
}

func (d *db) Tx(context context.Context) (pgx.Tx, error) {
	return d.pool.Begin(context)
}

func (d *db) Ping(context context.Context) error {
	return d.pool.Ping(context)
}
