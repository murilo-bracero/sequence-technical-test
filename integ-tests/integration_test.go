package integtests_test

import (
	"context"
	"testing"

	integtests "github.com/murilo-bracero/sequence-technical-test/integ-tests"
	"github.com/stretchr/testify/suite"
)

func TestIntegration(t *testing.T) {
	ev := integtests.New()
	ev.Start(context.Background())

	suite.Run(t, &StepHandlerTestSuite{ev: ev})
	suite.Run(t, &SequenceHandlerTestSuite{ev: ev})
}
