package services

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	dao "github.com/murilo-bracero/sequence-technical-test/internal/db/gen"
	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/murilo-bracero/sequence-technical-test/internal/repository"
)

type StepService interface {
	CreateStep(ctx context.Context, sequenceID uuid.UUID, req dto.CreateStepRequest) (*dto.StepResponse, error)
	UpdateStep(ctx context.Context, sequenceID uuid.UUID, stepID uuid.UUID, req dto.UpdateStepRequest) (*dto.StepResponse, error)
	DeleteStep(ctx context.Context, stepID uuid.UUID) error
}

type stepService struct {
	sequenceRepository repository.SequenceRepository
	stepRepository     repository.StepRepository
}

func NewStepService(sequenceRepository repository.SequenceRepository, stepRepository repository.StepRepository) StepService {
	return &stepService{sequenceRepository: sequenceRepository, stepRepository: stepRepository}
}

func (s *stepService) CreateStep(ctx context.Context, sequenceID uuid.UUID, req dto.CreateStepRequest) (*dto.StepResponse, error) {
	sequence, err := s.sequenceRepository.FindByExternalId(ctx, sequenceID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrorSequenceNotFound
		}
		slog.Error("failed to get sequence", err.Error(), err)
		return nil, err
	}

	step := &dao.Step{
		MailSubject: req.MailSubject,
		MailContent: req.MailContent,
		SequenceID:  sequence.ID,
	}

	if err := s.stepRepository.Create(ctx, step); err != nil {
		slog.Error("failed to create step", err.Error(), err)
		return nil, err
	}

	return &dto.StepResponse{
		ExternalID:  step.ExternalID.String(),
		MailSubject: step.MailSubject,
		MailContent: step.MailContent,
	}, nil
}

func (s *stepService) UpdateStep(ctx context.Context, sequenceID uuid.UUID, stepID uuid.UUID, req dto.UpdateStepRequest) (*dto.StepResponse, error) {
	step, err := s.stepRepository.FindOne(ctx, sequenceID, stepID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrorStepNotFound
		}

		slog.Error("failed to get step", err.Error(), err)
		return nil, err
	}

	if req.MailSubject != nil {
		step.MailSubject = *req.MailSubject
	}

	if req.MailContent != nil {
		step.MailContent = *req.MailContent
	}

	if err := s.stepRepository.Update(context.Background(), step); err != nil {
		slog.Error("failed to update step", err.Error(), err)
		return nil, err
	}

	return &dto.StepResponse{
		ExternalID:  step.ExternalID.String(),
		MailSubject: step.MailSubject,
		MailContent: step.MailContent,
	}, nil
}

func (s *stepService) DeleteStep(ctx context.Context, stepID uuid.UUID) error {
	return s.stepRepository.Delete(context.Background(), stepID)
}
