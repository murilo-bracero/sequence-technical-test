package services_test

import (
	"context"
	"database/sql"
	"testing"

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

func TestStepService_CreateStep(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("success", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		sequenceID := uuid.New()
		req := dto.CreateStepRequest{
			MailSubject: "subject",
			MailContent: "content",
		}

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(&models.SequenceWithSteps{ID: 1}, nil)
		stepRepository.EXPECT().Create(gomock.Any(), &dao.Step{
			MailSubject: req.MailSubject,
			MailContent: req.MailContent,
			SequenceID:  1,
		}).Return(nil)

		res, err := stepService.CreateStep(context.Background(), sequenceID, req)
		assert.NoError(t, err)
		assert.Equal(t, "subject", res.MailSubject)
		assert.Equal(t, "content", res.MailContent)
	})

	t.Run("return services.ErrorSequenceNotFound when sequence search fails with pgx.ErrNoRows", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		sequenceID := uuid.New()
		req := dto.CreateStepRequest{
			MailSubject: "subject",
			MailContent: "content",
		}

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(nil, pgx.ErrNoRows)
		stepRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		res, err := stepService.CreateStep(context.Background(), sequenceID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, services.ErrorSequenceNotFound.Error())
	})

	t.Run("return driver error in general cases", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		sequenceID := uuid.New()
		req := dto.CreateStepRequest{
			MailSubject: "subject",
			MailContent: "content",
		}

		sequenceRepository.EXPECT().FindByExternalId(gomock.Any(), sequenceID).Return(nil, sql.ErrConnDone)
		stepRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		res, err := stepService.CreateStep(context.Background(), sequenceID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})
}

func TestStepService_UpdateStep(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("success", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		sequenceID := uuid.New()
		stepID := uuid.New()

		mailSubject := "subject"
		mailContent := "content"

		req := dto.UpdateStepRequest{
			MailSubject: &mailSubject,
			MailContent: &mailContent,
		}

		stepRepository.EXPECT().FindOne(gomock.Any(), sequenceID, stepID).Return(&dao.Step{ID: 1}, nil)
		stepRepository.EXPECT().Update(gomock.Any(), &dao.Step{
			ID:          1,
			MailSubject: *req.MailSubject,
			MailContent: *req.MailContent,
		}).Return(nil)

		res, err := stepService.UpdateStep(context.Background(), sequenceID, stepID, req)
		assert.NoError(t, err)
		assert.Equal(t, "subject", res.MailSubject)
		assert.Equal(t, "content", res.MailContent)
	})

	t.Run("return ErrorStepNotFound when step search fails with pgx.ErrNoRows", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		sequenceID := uuid.New()
		stepID := uuid.New()

		mailSubject := "subject"
		mailContent := "content"

		req := dto.UpdateStepRequest{
			MailSubject: &mailSubject,
			MailContent: &mailContent,
		}

		stepRepository.EXPECT().FindOne(gomock.Any(), sequenceID, stepID).Return(nil, pgx.ErrNoRows)
		stepRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		res, err := stepService.UpdateStep(context.Background(), sequenceID, stepID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, services.ErrorStepNotFound.Error())
	})

	t.Run("return driver error in general cases", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		sequenceID := uuid.New()
		stepID := uuid.New()

		mailSubject := "subject"
		mailContent := "content"

		req := dto.UpdateStepRequest{
			MailSubject: &mailSubject,
			MailContent: &mailContent,
		}

		stepRepository.EXPECT().FindOne(gomock.Any(), sequenceID, stepID).Return(nil, sql.ErrConnDone)
		stepRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		res, err := stepService.UpdateStep(context.Background(), sequenceID, stepID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})

	t.Run("return driver error in general cases with second call", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		sequenceID := uuid.New()
		stepID := uuid.New()

		mailSubject := "subject"
		mailContent := "content"

		req := dto.UpdateStepRequest{
			MailSubject: &mailSubject,
			MailContent: &mailContent,
		}

		stepRepository.EXPECT().FindOne(gomock.Any(), sequenceID, stepID).Return(&dao.Step{ID: 1}, nil)
		stepRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(sql.ErrConnDone)

		res, err := stepService.UpdateStep(context.Background(), sequenceID, stepID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})
}

func TestStepService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("success", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		stepID := uuid.New()

		stepRepository.EXPECT().Delete(gomock.Any(), stepID).Return(nil)

		err := stepService.DeleteStep(context.Background(), stepID)
		assert.NoError(t, err)
	})

	t.Run("return driver error in general cases", func(t *testing.T) {
		sequenceRepository := mocks.NewMockSequenceRepository(ctrl)
		stepRepository := mocks.NewMockStepRepository(ctrl)
		stepService := services.NewStepService(sequenceRepository, stepRepository)

		stepID := uuid.New()

		stepRepository.EXPECT().Delete(gomock.Any(), stepID).Return(sql.ErrConnDone)

		err := stepService.DeleteStep(context.Background(), stepID)

		assert.EqualError(t, err, sql.ErrConnDone.Error())
	})
}
