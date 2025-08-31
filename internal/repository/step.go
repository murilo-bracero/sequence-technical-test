package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/murilo-bracero/sequence-technical-test/internal/db"
	dao "github.com/murilo-bracero/sequence-technical-test/internal/db/gen"
)

type StepRepository interface {
	FindOne(ctx context.Context, sequenceID uuid.UUID, stepID uuid.UUID) (*dao.Step, error)
	Create(ctx context.Context, model *dao.Step) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, model *dao.Step) error
}

type stepRepository struct {
	queries *dao.Queries
	db      db.DB
}

var _ StepRepository = (*stepRepository)(nil)

func NewStepRepository(db db.DB) *stepRepository {
	return &stepRepository{queries: db.Queries(), db: db}
}

func (r *stepRepository) FindOne(ctx context.Context, sequenceID uuid.UUID, stepID uuid.UUID) (*dao.Step, error) {
	step, err := r.queries.GetStepById(ctx, dao.GetStepByIdParams{ExternalID: stepID, ExternalID_2: sequenceID})
	if err != nil {
		return nil, err
	}

	return &step, nil
}

func (r *stepRepository) Create(ctx context.Context, model *dao.Step) error {
	step, err := r.queries.CreateStep(ctx, dao.CreateStepParams{
		SequenceID:  model.SequenceID,
		MailSubject: model.MailSubject,
		MailContent: model.MailContent,
	})
	if err != nil {
		return err
	}

	model.ID = step.ID
	model.ExternalID = step.ExternalID

	return nil
}

func (r *stepRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteStep(ctx, id)
}

func (r *stepRepository) Update(ctx context.Context, model *dao.Step) error {
	_, err := r.queries.UpdateStep(ctx, dao.UpdateStepParams{
		ExternalID:  model.ExternalID,
		MailSubject: model.MailSubject,
		MailContent: model.MailContent,
	})

	return err
}
