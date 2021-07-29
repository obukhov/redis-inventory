package renderer

import (
	"bytes"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JSONRendererTestSuite struct {
	suite.Suite
	trie *trie.Trie
}

func (suite *JSONRendererTestSuite) TestRender() {
	var buf bytes.Buffer

	r := JSONRenderer{&buf, JSONRendererParams{}}

	err := r.Render(suite.trie)
	suite.Assert().Nil(err, "Error rendering trie")

	suite.Assert().Equal(
		"{\"Children\":{\"dev:\":{\"Children\":{\"article:\":{\"Children\":{\"1\":{\"Values\":{\"Params\":{\"BytesSize\":100,\"KeysCount\":1}}},\"2\":{\"Values\":{\"Params\":{\"BytesSize\":100,\"KeysCount\":1}}}},\"Values\":{\"Params\":{\"BytesSize\":200,\"KeysCount\":2}}}}}},\"Values\":{\"Params\":{\"BytesSize\":200,\"KeysCount\":2}}}\n",
		buf.String(),
	)
}

func (suite *JSONRendererTestSuite) TestRenderWithIndent() {
	var buf bytes.Buffer

	params, err := NewJSONRendererParams("padSpaces=2")
	suite.Assert().Nil(err)

	r := JSONRenderer{&buf, params}

	err = r.Render(suite.trie)
	suite.Assert().Nil(err, "Error rendering trie")

	suite.Assert().Equal(
		"{\n"+
			"  \"Children\": {\n"+
			"    \"dev:\": {\n"+
			"      \"Children\": {\n"+
			"        \"article:\": {\n"+
			"          \"Children\": {\n"+
			"            \"1\": {\n"+
			"              \"Values\": {\n"+
			"                \"Params\": {\n"+
			"                  \"BytesSize\": 100,\n"+
			"                  \"KeysCount\": 1\n"+
			"                }\n"+
			"              }\n"+
			"            },\n"+
			"            \"2\": {\n"+
			"              \"Values\": {\n"+
			"                \"Params\": {\n"+
			"                  \"BytesSize\": 100,\n"+
			"                  \"KeysCount\": 1\n"+
			"                }\n"+
			"              }\n"+
			"            }\n"+
			"          },\n"+
			"          \"Values\": {\n"+
			"            \"Params\": {\n"+
			"              \"BytesSize\": 200,\n"+
			"              \"KeysCount\": 2\n"+
			"            }\n"+
			"          }\n"+
			"        }\n"+
			"      }\n"+
			"    }\n"+
			"  },\n"+
			"  \"Values\": {\n"+
			"    \"Params\": {\n"+
			"      \"BytesSize\": 200,\n"+
			"      \"KeysCount\": 2\n"+
			"    }\n"+
			"  }\n"+
			"}\n",
		buf.String(),
	)
}

func (suite *JSONRendererTestSuite) TestNewJSONRendererParams() {
	_, err := NewJSONRendererParams("padSpaces=asd")
	suite.Assert().Error(err)
}

func (suite *JSONRendererTestSuite) SetupTest() {
	suite.trie = trie.NewTrie(trie.NewPunctuationSplitter(':'), 3)

	suite.setupTrieKey("dev:article:1", 100)
	suite.setupTrieKey("dev:article:2", 100)
}

func (suite *JSONRendererTestSuite) setupTrieKey(key string, value int64) {
	suite.trie.Add(key, trie.ParamValue{Param: trie.BytesSize, Value: value}, trie.ParamValue{Param: trie.KeysCount, Value: 1})
}

func TestJsonRendererTestSuite(t *testing.T) {
	suite.Run(t, new(JSONRendererTestSuite))
}
