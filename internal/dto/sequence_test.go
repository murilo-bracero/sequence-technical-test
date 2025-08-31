package dto_test

import (
	"testing"

	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateSequenceRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		req := dto.CreateSequenceRequest{
			Name:                 "name",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: true,
			Steps: []*dto.CreateStepRequest{
				{
					StepNumber:  1,
					MailSubject: "subject",
					MailContent: "content",
				},
			},
		}
		assert.NoError(t, req.Validate())
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		req := dto.CreateSequenceRequest{
			Name:                 "",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: true,
			Steps: []*dto.CreateStepRequest{
				{
					StepNumber:  1,
					MailSubject: "subject",
					MailContent: "content",
				},
			},
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Equal(t, "sequence name is required", err.Error())
	})

	t.Run("should return error when open steps array is empty", func(t *testing.T) {
		req := dto.CreateSequenceRequest{
			Name:                 "name",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: true,
			Steps:                []*dto.CreateStepRequest{},
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Equal(t, "sequence steps are required", err.Error())
	})
}
