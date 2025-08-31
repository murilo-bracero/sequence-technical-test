package repository

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/murilo-bracero/sequence-technical-test/internal/db"
	dao "github.com/murilo-bracero/sequence-technical-test/internal/db/gen"
	"github.com/murilo-bracero/sequence-technical-test/internal/models"
)

type SequenceRepository interface {
	FindByExternalId(ctx context.Context, id uuid.UUID) (*models.SequenceWithSteps, error)
	FindAll(ctx context.Context, limit int, offset int) ([]*models.SequenceWithSteps, error)
	Create(ctx context.Context, model *models.SequenceWithSteps) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, model *models.SequenceWithSteps) error
}

type sequenceRepository struct {
	queries *dao.Queries
	db      db.DB
}

var _ SequenceRepository = (*sequenceRepository)(nil)

func NewSequenceRepository(db db.DB) *sequenceRepository {
	return &sequenceRepository{queries: db.Queries(), db: db}
}

func (r *sequenceRepository) FindByExternalId(ctx context.Context, externalId uuid.UUID) (*models.SequenceWithSteps, error) {
	row, err := r.queries.GetSequenceById(ctx, externalId)
	if err != nil {
		return nil, err
	}

	steps := make([]*dao.Step, 0)

	if err := json.Unmarshal(row.Steps, &steps); err != nil {
		slog.Error("failed to unmarshal steps", err.Error(), err)
	}

	model := &models.SequenceWithSteps{
		ID:                   row.ID,
		ExternalID:           row.ExternalID,
		Name:                 row.SequenceName,
		OpenTrackingEnabled:  row.OpenTrackingEnabled,
		ClickTrackingEnabled: row.ClickTrackingEnabled,
		Created:              row.Created.Time,
		Steps:                steps,
	}

	if row.Updated.Valid {
		model.Updated = &row.Updated.Time
	}

	return model, nil
}

func (r *sequenceRepository) FindAll(ctx context.Context, limit int, offset int) ([]*models.SequenceWithSteps, error) {
	rows, err := r.queries.GetSequences(ctx, dao.GetSequencesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	sequences := make([]*models.SequenceWithSteps, 0, len(rows))

	for _, row := range rows {
		steps := make([]*dao.Step, 0)

		if err := json.Unmarshal(row.Steps, &steps); err != nil {
			slog.Error("failed to unmarshal steps", err.Error(), err)
		}

		model := &models.SequenceWithSteps{
			ID:                   row.ID,
			ExternalID:           row.ExternalID,
			Name:                 row.SequenceName,
			OpenTrackingEnabled:  row.OpenTrackingEnabled,
			ClickTrackingEnabled: row.ClickTrackingEnabled,
			Created:              row.Created.Time,
			Steps:                steps,
		}

		if row.Updated.Valid {
			model.Updated = &row.Updated.Time
		}

		sequences = append(sequences, model)
	}

	return sequences, nil
}

func (r *sequenceRepository) Create(ctx context.Context, model *models.SequenceWithSteps) error {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", err.Error(), err)
		return err
	}

	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	sequence, err := qtx.CreateSequence(ctx, dao.CreateSequenceParams{
		SequenceName:         model.Name,
		OpenTrackingEnabled:  model.OpenTrackingEnabled,
		ClickTrackingEnabled: model.ClickTrackingEnabled,
	})

	model.ID = sequence.ID
	model.ExternalID = sequence.ExternalID
	model.Created = sequence.Created.Time

	if sequence.Updated.Valid {
		model.Updated = &sequence.Updated.Time
	}

	if err != nil {
		slog.Error("failed to create sequence", err.Error(), err)
		return err
	}

	createStepParams := make([]dao.CreateStepsParams, 0, len(model.Steps))
	for _, step := range model.Steps {
		step.ExternalID = uuid.New()
		step.SequenceID = sequence.ID

		createStepParams = append(createStepParams, dao.CreateStepsParams{
			ExternalID:  step.ExternalID,
			MailSubject: step.MailSubject,
			MailContent: step.MailContent,
			SequenceID:  step.SequenceID,
		})
	}

	_, err = qtx.CreateSteps(ctx, createStepParams)
	if err != nil {
		slog.Error("failed to create steps", err.Error(), err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("failed to commit transaction", err.Error(), err)
		return err
	}

	return nil
}

func (r *sequenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSequence(ctx, id)
}

func (r *sequenceRepository) Update(ctx context.Context, model *models.SequenceWithSteps) error {
	updated, err := r.queries.UpdateSequence(ctx, dao.UpdateSequenceParams{
		ID:                   model.ID,
		OpenTrackingEnabled:  model.OpenTrackingEnabled,
		ClickTrackingEnabled: model.ClickTrackingEnabled,
	})

	model.Updated = &updated.Updated.Time

	return err
}
