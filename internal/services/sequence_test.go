package services_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	dao "github.com/murilo-bracero/sequence-technical-test/internal/db/gen"
	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/murilo-bracero/sequence-technical-test/internal/models"
	"github.com/murilo-bracero/sequence-technical-test/internal/repository/mocks"
	"github.com/murilo-bracero/sequence-technical-test/internal/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSequeceService_GetSequences(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("success", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceRepository.EXPECT().FindAll(gomock.Any(), 10, 10).Return([]*models.SequenceWithSteps{
			{
				ID:                   1,
				ExternalID:           uuid.New(),
				Name:                 "name",
				OpenTrackingEnabled:  true,
				ClickTrackingEnabled: true,
				Steps: []*dao.Step{
					{
						ID:          1,
						ExternalID:  uuid.New(),
						MailSubject: "subject",
						MailContent: "content",
					},
				},
				Created: time.Now(),
				Updated: nil,
			},
		}, nil)

		res, err := sequenceService.GetSequences(context.Background(), 10, 1)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		assert.Equal(t, "name", res[0].Name)
		assert.NotEmpty(t, res[0].ExternalID)
		assert.NotNil(t, res[0].CreatedAt)
		assert.Nil(t, res[0].LastUpdatedAt)
		assert.NotEmpty(t, res[0].Steps[0].ExternalID)
		assert.Equal(t, "subject", res[0].Steps[0].MailSubject)
		assert.Equal(t, "content", res[0].Steps[0].MailContent)
	})

	t.Run("return general error in general cases", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceRepository.EXPECT().FindAll(gomock.Any(), 10, 10).Return(nil, sql.ErrConnDone)

		_, err := sequenceService.GetSequences(context.Background(), 10, 1)

		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})
}

func TestSequeceService_GetSequence(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("success", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceID := uuid.New()

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(&models.SequenceWithSteps{
			ID:                   1,
			ExternalID:           sequenceID,
			Name:                 "name",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: true,
			Steps: []*dao.Step{
				{
					ID:          1,
					ExternalID:  uuid.New(),
					MailSubject: "subject",
					MailContent: "content",
				},
			},
			Created: time.Now(),
			Updated: nil,
		}, nil)

		res, err := sequenceService.GetSequence(context.Background(), sequenceID)
		assert.NoError(t, err)

		assert.Equal(t, "name", res.Name)
		assert.NotEmpty(t, res.ExternalID)
		assert.NotNil(t, res.CreatedAt)
		assert.Nil(t, res.LastUpdatedAt)
		assert.NotEmpty(t, res.Steps[0].ExternalID)
		assert.Equal(t, "subject", res.Steps[0].MailSubject)
		assert.Equal(t, "content", res.Steps[0].MailContent)
	})

	t.Run("return services.ErrorSequenceNotFound when sequence search fails with pgx.ErrNoRows", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceID := uuid.New()

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(nil, pgx.ErrNoRows)

		_, err := sequenceService.GetSequence(context.Background(), sequenceID)

		assert.EqualError(t, err, services.ErrorSequenceNotFound.Error())
	})

	t.Run("return general error in general cases", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceID := uuid.New()

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(nil, sql.ErrConnDone)

		_, err := sequenceService.GetSequence(context.Background(), sequenceID)

		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})
}

func TestSequeceService_UpdateSequence(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("success", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceID := uuid.New()

		ope := true

		req := dto.UpdateSequenceRequest{
			OpenTrackingEnabled: &ope,
		}

		created := time.Now()

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(&models.SequenceWithSteps{
			ID:                   1,
			ExternalID:           sequenceID,
			Name:                 "name",
			OpenTrackingEnabled:  false,
			ClickTrackingEnabled: true,
			Steps:                nil,
			Created:              created,
			Updated:              nil,
		}, nil)

		sequenceRepository.EXPECT().Update(gomock.Any(), &models.SequenceWithSteps{
			ID:                   1,
			ExternalID:           sequenceID,
			Name:                 "name",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: true,
			Steps:                nil,
			Created:              created,
			Updated:              nil,
		}).Return(nil)

		res, err := sequenceService.UpdateSequence(context.Background(), sequenceID, req)
		assert.NoError(t, err)

		assert.Equal(t, "name", res.Name)
		assert.NotEmpty(t, res.ExternalID)
		assert.NotNil(t, res.CreatedAt)
		assert.Nil(t, res.LastUpdatedAt)
	})

	t.Run("return services.ErrorSequenceNotFound when sequence search fails with pgx.ErrNoRows", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceID := uuid.New()

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(nil, pgx.ErrNoRows)
		sequenceRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		_, err := sequenceService.UpdateSequence(context.Background(), sequenceID, dto.UpdateSequenceRequest{})

		assert.EqualError(t, err, services.ErrorSequenceNotFound.Error())
	})

	t.Run("return general error in general cases", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceID := uuid.New()

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(nil, sql.ErrConnDone)
		sequenceRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		_, err := sequenceService.UpdateSequence(context.Background(), sequenceID, dto.UpdateSequenceRequest{})

		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})

	t.Run("return general error in general cases when update", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceID := uuid.New()

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(&models.SequenceWithSteps{}, nil)
		sequenceRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(sql.ErrConnDone)

		_, err := sequenceService.UpdateSequence(context.Background(), sequenceID, dto.UpdateSequenceRequest{})

		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})
}

func TestSequeceService_CreateSequence(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("success", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		step := &dto.CreateStepRequest{
			MailSubject: "subject",
			MailContent: "content",
		}

		req := dto.CreateSequenceRequest{
			Name:                 "name",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: true,
			Steps:                []*dto.CreateStepRequest{step},
		}

		sequenceRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		res, err := sequenceService.CreateSequence(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "name", res.Name)
		assert.NotEmpty(t, res.ExternalID)
		assert.NotNil(t, res.CreatedAt)
		assert.Nil(t, res.LastUpdatedAt)

		assert.Len(t, res.Steps, 1)
		assert.Equal(t, "subject", res.Steps[0].MailSubject)
		assert.Equal(t, "content", res.Steps[0].MailContent)
	})

	t.Run("return general error in general cases", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		sequenceService := services.NewSequenceService(sequenceRepository)

		sequenceRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(sql.ErrConnDone)

		_, err := sequenceService.CreateSequence(context.Background(), dto.CreateSequenceRequest{})

		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})
}
