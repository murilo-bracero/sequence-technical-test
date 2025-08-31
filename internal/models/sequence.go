package models

import (
	"time"

	"github.com/google/uuid"
	dao "github.com/murilo-bracero/sequence-technical-test/internal/db/gen"
)

type SequenceWithSteps struct {
	ID                   int32
	ExternalID           uuid.UUID
	Name                 string
	OpenTrackingEnabled  bool
	ClickTrackingEnabled bool
	Created              time.Time
	Updated              *time.Time
	Steps                []*dao.Step
}
