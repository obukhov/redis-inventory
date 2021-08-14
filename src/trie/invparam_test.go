package trie

import (
	"encoding/json"
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

func (suite *InvParamTestSuite) TestUnmarshalText() {
	jsonString := `["BytesSize", "KeysCount"]`
	var result []InvParam

	err := json.Unmarshal([]byte(jsonString), &result)

	suite.Assert().Nil(err)
	suite.Assert().Equal([]InvParam{BytesSize, KeysCount}, result)
}

func (suite *InvParamTestSuite) TestUnmarshalTextError() {
	jsonString := `["BytesSizeFoo", "KeysCount"]`
	var result []InvParam

	err := json.Unmarshal([]byte(jsonString), &result)

	suite.Assert().NotNil(err)
}

func TestInvParamTestSuite(t *testing.T) {
	suite.Run(t, new(InvParamTestSuite))
}
