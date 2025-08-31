package integtests_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	integtests "github.com/murilo-bracero/sequence-technical-test/integ-tests"
	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type SequenceHandlerTestSuite struct {
	suite.Suite
	ev *integtests.EnvironmentCommands
}

func (s *SequenceHandlerTestSuite) SetupSuite() {
	if s.ev == nil {
		s.T().Fatal("No environment created")
	}
}

func (s *SequenceHandlerTestSuite) TestSequenceHandler_CreateSequence() {
	t := s.T()

	url := "http://localhost:8000/sequences"

	payload := strings.NewReader("{\"name\": \"My Sequence 982\",\"openTrackingEnabled\": false,\"clickTrackingEnabled\": true,\"steps\": [{\"mailSubject\": \"Subject 96\",\"mailContent\": \"78 Sat, 30 Aug 2025 23:51:34 GMT\"}]}")

	req, err := http.NewRequest("POST", url, payload)

	assert.NoError(t, err)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, res.StatusCode)

	var body dto.SequenceResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "My Sequence 982", body.Name)
	assert.Equal(t, false, body.OpenTrackingEnabled)
	assert.Equal(t, true, body.ClickTrackingEnabled)
	assert.Len(t, body.Steps, 1)
}

func (s *SequenceHandlerTestSuite) TestSequenceHandler_CreateSequence_BadRequest() {
	t := s.T()

	url := "http://localhost:8000/sequences"

	payload := strings.NewReader("{\"name\": \"\",\"openTrackingEnabled\": false,\"clickTrackingEnabled\": true,\"steps\": [{\"mailSubject\": \"Subject 96\",\"mailContent\": \"78 Sat, 30 Aug 2025 23:51:34 GMT\"}]}")

	req, err := http.NewRequest("POST", url, payload)

	assert.NoError(t, err)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

func (s *SequenceHandlerTestSuite) TestSequenceHandler_GetSequences() {
	t := s.T()

	url := "http://localhost:8000/sequences"

	req, err := http.NewRequest("GET", url, nil)

	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	var body []dto.SequenceResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, body, 1)

	assert.Equal(t, "My Sequence 982", body[0].Name)
	assert.Equal(t, false, body[0].OpenTrackingEnabled)
	assert.Equal(t, true, body[0].ClickTrackingEnabled)
	assert.Len(t, body[0].Steps, 1)
}

func (s *SequenceHandlerTestSuite) TestSequenceHandler_GetSequence() {
	t := s.T()

	url := "http://localhost:8000/sequences"

	req, err := http.NewRequest("GET", url, nil)

	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	var body []dto.SequenceResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, body, 1)

	id := body[0].ExternalID

	url = "http://localhost:8000/sequences/" + id

	req, err = http.NewRequest("GET", url, nil)

	assert.NoError(t, err)

	res, err = http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	var body2 dto.SequenceResponse
	if err := json.NewDecoder(res.Body).Decode(&body2); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "My Sequence 982", body2.Name)
	assert.Equal(t, false, body2.OpenTrackingEnabled)
	assert.Equal(t, true, body2.ClickTrackingEnabled)
	assert.Len(t, body2.Steps, 1)
}

func (s *SequenceHandlerTestSuite) TestSequenceHandler_UpdateSequence() {
	t := s.T()

	url := "http://localhost:8000/sequences"

	req, err := http.NewRequest("GET", url, nil)

	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	var listResponse []dto.SequenceResponse
	if err := json.NewDecoder(res.Body).Decode(&listResponse); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, listResponse, 1)

	id := listResponse[0].ExternalID

	url = "http://localhost:8000/sequences/" + id

	payload := strings.NewReader("{\"openTrackingEnabled\": true,\"clickTrackingEnabled\": false}")

	req, err = http.NewRequest("PATCH", url, payload)

	assert.NoError(t, err)

	req.Header.Add("content-type", "application/json")

	res, err = http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	var patchResponse dto.SequenceResponse
	if err := json.NewDecoder(res.Body).Decode(&patchResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "My Sequence 982", patchResponse.Name)
	assert.Equal(t, true, patchResponse.OpenTrackingEnabled)
	assert.Equal(t, false, patchResponse.ClickTrackingEnabled)
	assert.NotEmpty(t, patchResponse.LastUpdatedAt)
	assert.Len(t, patchResponse.Steps, 1)
}

func (s *SequenceHandlerTestSuite) TearDownSuite() {
	err := s.ev.ClearDatabase(context.Background())
	s.Require().NoError(err)
}
