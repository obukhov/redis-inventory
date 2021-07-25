package renderer

import (
	"bytes"
	"fmt"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TableRendererTestSuite struct {
	suite.Suite
	trie *trie.Trie
}

func (suite *TableRendererTestSuite) TestRenderSpacePadding() {
	var buf bytes.Buffer

	r := TableRenderer{&buf, TableRendererParams{10, "", 2, false, "  "}}

	err := r.Render(suite.trie)
	suite.Assert().Nil(err, "Error rendering trie")

	suite.Assert().Equal(
		""+
			"+-----------------------+----------+-----------+\n"+
			"| KEY                   | BYTESIZE | KEYSCOUNT |\n"+
			"+-----------------------+----------+-----------+\n"+
			"|   dev:                |     2500 |         7 |\n"+
			"|     article:          |      500 |         5 |\n"+
			"|       1               |      100 |         1 |\n"+
			"|       2               |      100 |         1 |\n"+
			"|       3               |      100 |         1 |\n"+
			"|       ( 2 more keys ) |          |           |\n"+
			"|     user:             |     2000 |         2 |\n"+
			"|       bar             |     1000 |         1 |\n"+
			"|       foo             |     1000 |         1 |\n"+
			"|   prod:user:          |     4000 |         2 |\n"+
			"|       bar             |     2000 |         1 |\n"+
			"|       foo             |     2000 |         1 |\n"+
			"+-----------------------+----------+-----------+\n",
		buf.String(),
	)
}

func (suite *TableRendererTestSuite) TestRenderFullPathAndDepthLimit() {
	var buf bytes.Buffer

	r := TableRenderer{&buf, TableRendererParams{Depth: 2}}

	err := r.Render(suite.trie)
	suite.Assert().Nil(err, "Error rendering trie")

	suite.Assert().Equal(
		""+
			"+---------------+----------+-----------+\n"+
			"| KEY           | BYTESIZE | KEYSCOUNT |\n"+
			"+---------------+----------+-----------+\n"+
			"| dev:          |     2500 |         7 |\n"+
			"| dev:article:  |      500 |         5 |\n"+
			"| dev:user:     |     2000 |         2 |\n"+
			"| prod:user:    |     4000 |         2 |\n"+
			"| prod:user:bar |     2000 |         1 |\n"+
			"| prod:user:foo |     2000 |         1 |\n"+
			"+---------------+----------+-----------+\n",
		buf.String(),
	)
}

func (suite *TableRendererTestSuite) TestRenderHuman() {
	var buf bytes.Buffer

	r := TableRenderer{&buf, TableRendererParams{Depth: 2, HumanReadable: true}}

	for i := 0; i < 1000; i++ {
		suite.setupTrieKey(fmt.Sprintf("dev:blog:%d", i), 2)
	}

	err := r.Render(suite.trie)
	suite.Assert().Nil(err, "Error rendering trie")

	suite.Assert().Equal(
		""+
			"+---------------+----------+-----------+\n"+
			"| KEY           | BYTESIZE | KEYSCOUNT |\n"+
			"+---------------+----------+-----------+\n"+
			"| dev:          |     4.4K |     1,007 |\n"+
			"| dev:article:  |     500B |         5 |\n"+
			"| dev:blog:     |       2K |     1,000 |\n"+
			"| dev:user:     |       2K |         2 |\n"+
			"| prod:user:    |     3.9K |         2 |\n"+
			"| prod:user:bar |       2K |         1 |\n"+
			"| prod:user:foo |       2K |         1 |\n"+
			"+---------------+----------+-----------+\n",
		buf.String(),
	)
}

func (suite *TableRendererTestSuite) TestRenderParam() {
	for _, t := range []struct {
		paramString    string
		expectedParams TableRendererParams
	}{
		{
			"",
			TableRendererParams{Depth: 10}},
		{
			"depth=2&padSpaces=3",
			TableRendererParams{Depth: 2, PaddingSpaceCount: 3, indent: "   "},
		},
		{
			"depth=5&padding=>>",
			TableRendererParams{Depth: 5, Padding: ">>", indent: ">>"},
		},
	} {
		suite.Run(t.paramString, func() {

			params, err := NewTableRendererParams(t.paramString)
			suite.Assert().Nil(err)
			suite.Assert().Equal(t.expectedParams, params)
		})
	}
}

func (suite *TableRendererTestSuite) SetupTest() {
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

func (suite *TableRendererTestSuite) setupTrieKey(key string, value int64) {
	suite.trie.Add(key, trie.ParamValue{Param: trie.BytesSize, Value: value}, trie.ParamValue{Param: trie.KeysCount, Value: 1})
}

func TestTableRendererTestSuite(t *testing.T) {
	suite.Run(t, new(TableRendererTestSuite))
}
