package dto

type UpdateStepRequest struct {
	MailSubject *string `json:"mailSubject"`
	MailContent *string `json:"mailContent"`
}
