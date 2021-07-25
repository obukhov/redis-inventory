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
		"{\"Children\":{\"dev:\":{\"Children\":{\"article:\":{\"Children\":{\"1\":{\"Values\":{\"Params\":{\"BytesSize\":100,\"KeysCount\":1}}},\"2\":{\"Values\":{\"Params\":{\"BytesSize\":100,\"KeysCount\":1}}},\"3\":{\"Values\":{\"Params\":{\"BytesSize\":100,\"KeysCount\":1}}}},\"Values\":{\"Params\":{\"BytesSize\":500,\"KeysCount\":5}},\"Overflow\":2},\"user:\":{\"Children\":{\"bar\":{\"Values\":{\"Params\":{\"BytesSize\":1000,\"KeysCount\":1}}},\"foo\":{\"Values\":{\"Params\":{\"BytesSize\":1000,\"KeysCount\":1}}}},\"Values\":{\"Params\":{\"BytesSize\":2000,\"KeysCount\":2}}}},\"Values\":{\"Params\":{\"BytesSize\":2500,\"KeysCount\":7}}},\"prod:\":{\"Children\":{\"user:\":{\"Children\":{\"bar\":{\"Values\":{\"Params\":{\"BytesSize\":2000,\"KeysCount\":1}}},\"foo\":{\"Values\":{\"Params\":{\"BytesSize\":2000,\"KeysCount\":1}}}},\"Values\":{\"Params\":{\"BytesSize\":4000,\"KeysCount\":2}}}}}},\"Values\":{\"Params\":{\"BytesSize\":6500,\"KeysCount\":9}}}\n",
		buf.String(),
	)
}

func (suite *JSONRendererTestSuite) SetupTest() {
	suite.trie = trie.NewTrie(trie.NewPunctuationSplitter(':'), 3)

	suite.setupTrieKey("dev:article:1", 100)
	suite.setupTrieKey("dev:article:2", 100)
	suite.setupTrieKey("dev:article:3", 100)
	suite.setupTrieKey("dev:article:4", 100)
	suite.setupTrieKey("dev:article:5", 100)
	suite.setupTrieKey("dev:user:bar", 1000)
	suite.setupTrieKey("dev:user:foo", 1000)
	suite.setupTrieKey("prod:user:bar", 2000)
	suite.setupTrieKey("prod:user:foo", 2000)
}

func (suite *JSONRendererTestSuite) setupTrieKey(key string, value int64) {
	suite.trie.Add(key, trie.ParamValue{Param: trie.BytesSize, Value: value}, trie.ParamValue{Param: trie.KeysCount, Value: 1})
}

func TestJsonRendererTestSuite(t *testing.T) {
	suite.Run(t, new(JSONRendererTestSuite))
}
