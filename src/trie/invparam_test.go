package trie

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type InvParamTestSuite struct {
	suite.Suite
}

func (suite *InvParamTestSuite) TestStringPanic() {
	suite.Assert().Panics(func() {
		p := InvParam(9999)
		_ = p.String()
	})
}

func TestInvParamTestSuite(t *testing.T) {
	suite.Run(t, new(InvParamTestSuite))
}
