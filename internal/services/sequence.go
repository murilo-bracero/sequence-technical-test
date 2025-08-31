package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	dao "github.com/murilo-bracero/sequence-technical-test/internal/db/gen"
	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/murilo-bracero/sequence-technical-test/internal/models"
	"github.com/murilo-bracero/sequence-technical-test/internal/repository"
)

type SequenceService interface {
	GetSequences(ctx context.Context, size int, page int) ([]*dto.SequenceResponse, error)
	GetSequence(ctx context.Context, id uuid.UUID) (*dto.SequenceResponse, error)
	UpdateSequence(ctx context.Context, id uuid.UUID, req dto.UpdateSequenceRequest) (*dto.SequenceResponse, error)
	CreateSequence(ctx context.Context, req dto.CreateSequenceRequest) (*dto.SequenceResponse, error)
}

type sequenceService struct {
	sequenceRepository repository.SequenceRepository
}

func NewSequenceService(sequenceRepository repository.SequenceRepository) SequenceService {
	return &sequenceService{sequenceRepository: sequenceRepository}
}

func (s *sequenceService) GetSequences(ctx context.Context, size int, page int) ([]*dto.SequenceResponse, error) {
	sequences, err := s.sequenceRepository.FindAll(ctx, size, size*page)
	if err != nil {
		slog.Error("failed to get sequences", err.Error(), err)
		return nil, err
	}

	response := make([]*dto.SequenceResponse, 0, len(sequences))
	for _, s := range sequences {
		sr := &dto.SequenceResponse{
			ExternalID:           s.ExternalID.String(),
			Name:                 s.Name,
			OpenTrackingEnabled:  s.OpenTrackingEnabled,
			ClickTrackingEnabled: s.ClickTrackingEnabled,
			CreatedAt:            s.Created.Format(time.RFC3339),
			Steps:                make([]*dto.StepResponse, 0, len(s.Steps)),
		}

		if s.Updated != nil {
			updated := s.Updated.Format(time.RFC3339)
			sr.LastUpdatedAt = &updated
		}

		for _, step := range s.Steps {
			if step == nil {
				continue
			}

			sr.Steps = append(sr.Steps, &dto.StepResponse{
				ExternalID:  step.ExternalID.String(),
				MailSubject: step.MailSubject,
				MailContent: step.MailContent,
			})
		}

		response = append(response, sr)
	}

	return response, nil
}

func (s *sequenceService) GetSequence(ctx context.Context, id uuid.UUID) (*dto.SequenceResponse, error) {
	sequence, err := s.sequenceRepository.FindByExternalId(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrorSequenceNotFound
		}

		slog.Error("failed to get sequence", err.Error(), err)
		return nil, err
	}

	response := &dto.SequenceResponse{
		ExternalID:           sequence.ExternalID.String(),
		Name:                 sequence.Name,
		OpenTrackingEnabled:  sequence.OpenTrackingEnabled,
		ClickTrackingEnabled: sequence.ClickTrackingEnabled,
		CreatedAt:            sequence.Created.Format(time.RFC3339),
		Steps:                make([]*dto.StepResponse, 0, len(sequence.Steps)),
	}

	if sequence.Updated != nil {
		updated := sequence.Updated.Format(time.RFC3339)
		response.LastUpdatedAt = &updated
	}

	for _, step := range sequence.Steps {
		if step == nil {
			continue
		}
		response.Steps = append(response.Steps, &dto.StepResponse{
			ExternalID:  step.ExternalID.String(),
			MailSubject: step.MailSubject,
			MailContent: step.MailContent,
		})
	}

	return response, nil
}

func (s *sequenceService) UpdateSequence(ctx context.Context, id uuid.UUID, req dto.UpdateSequenceRequest) (*dto.SequenceResponse, error) {
	sequence, err := s.sequenceRepository.FindByExternalId(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrorSequenceNotFound
		}
		slog.Error("failed to get sequence during updateSequence", err.Error(), err)
		return nil, err
	}

	if req.OpenTrackingEnabled != nil {
		sequence.OpenTrackingEnabled = *req.OpenTrackingEnabled
	}

	if req.ClickTrackingEnabled != nil {
		sequence.ClickTrackingEnabled = *req.ClickTrackingEnabled
	}

	if err := s.sequenceRepository.Update(ctx, sequence); err != nil {
		slog.Error("failed to update sequence", err.Error(), err)
		return nil, err
	}

	response := &dto.SequenceResponse{
		ExternalID:           sequence.ExternalID.String(),
		Name:                 sequence.Name,
		OpenTrackingEnabled:  sequence.OpenTrackingEnabled,
		ClickTrackingEnabled: sequence.ClickTrackingEnabled,
		CreatedAt:            sequence.Created.Format(time.RFC3339),
	}

	if sequence.Updated != nil {
		updated := sequence.Updated.Format(time.RFC3339)
		response.LastUpdatedAt = &updated
	}

	for _, step := range sequence.Steps {
		response.Steps = append(response.Steps, &dto.StepResponse{
			ExternalID:  step.ExternalID.String(),
			MailSubject: step.MailSubject,
			MailContent: step.MailContent,
		})
	}

	return response, nil
}

func (s *sequenceService) CreateSequence(ctx context.Context, req dto.CreateSequenceRequest) (*dto.SequenceResponse, error) {
	sequence := models.SequenceWithSteps{
		Name:                 req.Name,
		OpenTrackingEnabled:  req.OpenTrackingEnabled,
		ClickTrackingEnabled: req.ClickTrackingEnabled,
		Steps:                make([]*dao.Step, 0, len(req.Steps)),
	}

	for _, step := range req.Steps {
		sequence.Steps = append(sequence.Steps, &dao.Step{
			MailSubject: step.MailSubject,
			MailContent: step.MailContent,
		})
	}

	if err := s.sequenceRepository.Create(ctx, &sequence); err != nil {
		slog.Error("failed to create sequence", err.Error(), err)
		return nil, err
	}

	response := &dto.SequenceResponse{
		ExternalID:           sequence.ExternalID.String(),
		Name:                 sequence.Name,
		OpenTrackingEnabled:  sequence.OpenTrackingEnabled,
		ClickTrackingEnabled: sequence.ClickTrackingEnabled,
		CreatedAt:            sequence.Created.Format(time.RFC3339),
		Steps:                make([]*dto.StepResponse, 0, len(req.Steps)),
	}

	for _, step := range sequence.Steps {
		response.Steps = append(response.Steps, &dto.StepResponse{
			ExternalID:  step.ExternalID.String(),
			MailSubject: step.MailSubject,
			MailContent: step.MailContent,
		})
	}

	return response, nil
}
