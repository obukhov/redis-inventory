package trie

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type InvParamTestSuite struct {
	suite.Suite
}


func (suite *InvParamTestSuite) TestKeyPanic() {
	suite.Assert().Panics(suite.T(), func(){
		p := InvParam(9999)
		p.String()
	})
}

func TestInvParamTestSuite(t *testing.T) {
	suite.Run(t, new(InvParamTestSuite))
}
