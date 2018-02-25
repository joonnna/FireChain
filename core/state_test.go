package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StateTestSuite struct {
	suite.Suite
	p *peer
}

func TestStateTestSuite(t *testing.T) {
	suite.Run(t, new(StateTestSuite))
}

func (suite *StateTestSuite) SetupTest() {
	suite.p = &peer{}
}
