package dto_test

import (
	"testing"

	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateStepRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		req := dto.CreateStepRequest{
			MailSubject: "subject",
			MailContent: "content",
		}
		assert.NoError(t, req.Validate())
	})

	t.Run("should return error when mail subject is empty", func(t *testing.T) {
		req := dto.CreateStepRequest{
			MailSubject: "",
			MailContent: "content",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Equal(t, "mail subject is required", err.Error())
	})

	t.Run("should return error when mail content is empty", func(t *testing.T) {
		req := dto.CreateStepRequest{
			MailSubject: "subject",
			MailContent: "",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Equal(t, "mail content is required", err.Error())
	})
}
