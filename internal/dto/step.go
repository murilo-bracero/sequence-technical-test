package dto

import "fmt"

type UpdateStepRequest struct {
	MailSubject *string `json:"mailSubject"`
	MailContent *string `json:"mailContent"`
}

type CreateStepRequest struct {
	MailSubject string `json:"mailSubject"`
	MailContent string `json:"mailContent"`
}

func (req *CreateStepRequest) Validate() error {
	if req.MailSubject == "" {
		return fmt.Errorf("mail subject is required")
	}
	if req.MailContent == "" {
		return fmt.Errorf("mail content is required")
	}
	return nil
}
