package dto

import "fmt"

type CreateSequenceRequest struct {
	Name                 string               `json:"Name"`
	OpenTrackingEnabled  bool                 `json:"openTrackingEnabled"`
	ClickTrackingEnabled bool                 `json:"clickTrackingEnabled"`
	Steps                []*CreateStepRequest `json:"steps"`
}

func (req *CreateSequenceRequest) Validate() error {
	if req.Name == "" {
		return fmt.Errorf("sequence name is required")
	}
	if len(req.Steps) == 0 {
		return fmt.Errorf("sequence steps are required")
	}

	for _, step := range req.Steps {
		if err := step.Validate(); err != nil {
			return err
		}
	}

	// checks if the step numbers are unique
	stepNumbers := make(map[int]bool)
	for _, step := range req.Steps {
		if _, ok := stepNumbers[step.StepNumber]; ok {
			return fmt.Errorf("step number %d is not unique", step.StepNumber)
		}
		stepNumbers[step.StepNumber] = true
	}

	return nil
}

type UpdateSequenceRequest struct {
	OpenTrackingEnabled  *bool `json:"openTrackingEnabled"`
	ClickTrackingEnabled *bool `json:"clickTrackingEnabled"`
}

type SequenceResponse struct {
	ExternalID           string          `json:"id"`
	Name                 string          `json:"name"`
	OpenTrackingEnabled  bool            `json:"openTrackingEnabled"`
	ClickTrackingEnabled bool            `json:"clickTrackingEnabled"`
	Steps                []*StepResponse `json:"steps"`
	CreatedAt            string          `json:"createdAt"`
	LastUpdatedAt        *string         `json:"lastUpdatedAt"`
}

type StepResponse struct {
	ExternalID  string `json:"id"`
	StepNumber  int    `json:"stepNumber"`
	MailSubject string `json:"mailSubject"`
	MailContent string `json:"mailContent"`
}
