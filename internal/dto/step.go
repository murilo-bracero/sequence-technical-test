package dto

import "fmt"

type UpdateStepRequest struct {
	StepNumber  *int    `json:"stepNumber"`
	MailSubject *string `json:"mailSubject"`
	MailContent *string `json:"mailContent"`
}

type CreateStepRequest struct {
	StepNumber  int    `json:"stepNumber"`
	MailSubject string `json:"mailSubject"`
	MailContent string `json:"mailContent"`
}

func (req *CreateStepRequest) Validate() error {
	if req.StepNumber == 0 {
		return fmt.Errorf("step number is required")
	}

	if req.MailSubject == "" {
		return fmt.Errorf("mail subject is required")
	}
	if req.MailContent == "" {
		return fmt.Errorf("mail content is required")
	}
	return nil
}
