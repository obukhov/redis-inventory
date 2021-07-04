package trie

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TrieTestSuite struct {
	suite.Suite
	splitter Splitter
}

func (suite *TrieTestSuite) TestAdd() {
	trie := NewTrie(NewPunctuationSplitter(':'), 100)

	trie.Add("foo:bar:lorem", ParamValue{BytesSize, 10})
	trie.Add("foo:bar:ipsum", ParamValue{BytesSize, 20})

	assert.Equal(suite.T(), int64(30), trie.root.Aggr.Params[BytesSize], "Root node aggregated value is 30")

	assert.Nil(suite.T(), trie.root.Children["foo:"].Aggr, "Intermediate node skip aggregation if has just one child")

	assert.NotNil(suite.T(), trie.root.Children["foo:"].Children["bar:"].Aggr, "Fork node has aggregator")
	assert.Equal(suite.T(), int64(30), trie.root.Children["foo"].Children["bar"].Aggr.Params[BytesSize], "Fork node aggregated value is 30")
}

func TestTrieTestSuite(t *testing.T) {
	suite.Run(t, new(TrieTestSuite))
}
