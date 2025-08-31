package dto

type CreateSequenceRequest struct {
	Name                 string               `json:"Name"`
	OpenTrackingEnabled  bool                 `json:"openTrackingEnabled"`
	ClickTrackingEnabled bool                 `json:"clickTrackingEnabled"`
	Steps                []*CreateStepRequest `json:"steps"`
}

type CreateStepRequest struct {
	MailSubject string `json:"mailSubject"`
	MailContent string `json:"mailContent"`
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
	MailSubject string `json:"mailSubject"`
	MailContent string `json:"mailContent"`
}
