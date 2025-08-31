package integtests_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	integtests "github.com/murilo-bracero/sequence-technical-test/integ-tests"
	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type StepHandlerTestSuite struct {
	suite.Suite
	ev *integtests.EnvironmentCommands
}

func (s *StepHandlerTestSuite) SetupSuite() {
	if s.ev == nil {
		s.T().Fatal("No environment created")
	}
}

func (s *StepHandlerTestSuite) TestStepHandler_CreateStep() {
	t := s.T()

	sequence, err := s.ev.CreateSequence(context.Background(), dto.CreateSequenceRequest{
		Name:                 "My Sequence 1",
		OpenTrackingEnabled:  false,
		ClickTrackingEnabled: true,
		Steps:                []*dto.CreateStepRequest{{MailSubject: "test subject", MailContent: "test mailbody", StepNumber: 1}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, sequence)

	url := fmt.Sprintf("http://localhost:8000/sequences/%s/steps", sequence.ExternalID)

	payload := strings.NewReader(`
	{
	"stepNumber": 3,
    "mailSubject": "test subject",
    "mailContent": "test mailbody"
	}
	`)

	req, err := http.NewRequest("POST", url, payload)

	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var body dto.StepResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, body.ExternalID)

	assert.Equal(t, "test subject", body.MailSubject)
	assert.Equal(t, "test mailbody", body.MailContent)
}

func (s *StepHandlerTestSuite) TestStepHandler_UpdateStep() {
	t := s.T()

	sequence, err := s.ev.CreateSequence(context.Background(), dto.CreateSequenceRequest{
		Name:                 "My Sequence 1",
		OpenTrackingEnabled:  false,
		ClickTrackingEnabled: true,
		Steps:                []*dto.CreateStepRequest{{MailSubject: "test subject", MailContent: "test mailbody", StepNumber: 1}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, sequence)

	url := fmt.Sprintf("http://localhost:8000/sequences/%s/steps/%s", sequence.ExternalID, sequence.Steps[0].ExternalID)

	payload := strings.NewReader(`
	{
    "mailSubject": "test subject",
    "mailContent": "test mailbody"
	}
	`)

	req, err := http.NewRequest("PATCH", url, payload)

	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var body dto.StepResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, body.ExternalID)

	assert.Equal(t, "test subject", body.MailSubject)
	assert.Equal(t, "test mailbody", body.MailContent)
}

func (s *StepHandlerTestSuite) TestStepHandler_DeleteStep() {
	t := s.T()

	sequence, err := s.ev.CreateSequence(context.Background(), dto.CreateSequenceRequest{
		Name:                 "My Sequence 1",
		OpenTrackingEnabled:  false,
		ClickTrackingEnabled: true,
		Steps:                []*dto.CreateStepRequest{{MailSubject: "test subject", MailContent: "test mailbody", StepNumber: 1}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, sequence)

	url := fmt.Sprintf("http://localhost:8000/sequences/%s/steps/%s", sequence.ExternalID, sequence.Steps[0].ExternalID)

	req, err := http.NewRequest("DELETE", url, nil)

	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	sequence, err = s.ev.GetSequenceById(context.Background(), sequence.ExternalID)

	assert.NoError(t, err)
	assert.NotNil(t, sequence)
	assert.Len(t, sequence.Steps, 0)
}

func (s *StepHandlerTestSuite) TearDownSuite() {
	err := s.ev.ClearDatabase(context.Background())
	s.Require().NoError(err)
}
