package services

import "errors"

var (
	ErrorSequenceNotFound = errors.New("sequence not found")
	ErrorStepNotFound     = errors.New("step not found")
)
