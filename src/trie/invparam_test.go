package trie

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InvParamTestSuite struct {
	suite.Suite
}


func (suite *InvParamTestSuite) TestKeyPanic() {
  suite.Assert().Panics(t, func(){
    p := InvParam(9999)
    p.String()
  })
}

func TestInvParamTestSuite(t *testing.T) {
	suite.Run(t, new(InvParamTestSuite))
}
