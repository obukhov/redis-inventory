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
	// root node
	assert.Equal(suite.T(), int64(30), trie.root.Aggr.Params[BytesSize], "Root node aggregated value is 30")
	// intermediate node
	assert.Nil(suite.T(), trie.root.Children["foo:"].Aggr, "Intermediate node skip aggregation if has just one child")
	// fork node
	assert.NotNil(suite.T(), trie.root.Children["foo:"].Children["bar:"].Aggr, "Fork node has aggregator")
	assert.Equal(suite.T(), int64(30), trie.root.Children["foo:"].Children["bar:"].Aggr.Params[BytesSize], "Fork node aggregated value is 30")
}

func (suite *TrieTestSuite) TestKeyEndingInFork() {
	trie := NewTrie(NewPunctuationSplitter(':'), 100)

	trie.Add("foo:bar:lorem", ParamValue{BytesSize, 10})
	trie.Add("foo:bar:ipsum", ParamValue{BytesSize, 20})
	trie.Add("foo:bar:", ParamValue{BytesSize, 100})
	// root node
	assert.Equal(suite.T(), int64(130), trie.root.Aggr.Params[BytesSize], "Root node aggregated value is 130")
	// intermediate node
	assert.Nil(suite.T(), trie.root.Children["foo:"].Aggr, "Intermediate node skip aggregation if has just one child")
	// fork node
	assert.NotNil(suite.T(), trie.root.Children["foo:"].Children["bar:"].Aggr, "Fork node has where key ended")
	assert.Equal(suite.T(), int64(130), trie.root.Children["foo:"].Children["bar:"].Aggr.Params[BytesSize], "Fork node aggregated value is 130")
}

func (suite *TrieTestSuite) TestKeyEndingInIntermediateNode() {
	trie := NewTrie(NewPunctuationSplitter(':'), 100)

	trie.Add("foo:bar:lorem", ParamValue{BytesSize, 10})
	trie.Add("foo:bar:ipsum", ParamValue{BytesSize, 20})
	trie.Add("foo:", ParamValue{BytesSize, 100})
	// root node
	assert.Equal(suite.T(), int64(130), trie.root.Aggr.Params[BytesSize], "Root node aggregated value is 30")
	// intermediate node
	assert.NotNil(suite.T(), trie.root.Children["foo:"].Aggr, "Intermediate node has aggregator if key ends here")
	assert.Equal(suite.T(), int64(130), trie.root.Children["foo:"].Aggr.Params[BytesSize], "Fork node aggregated value is 130")
	// fork node
	assert.NotNil(suite.T(), trie.root.Children["foo:"].Children["bar:"].Aggr, "Fork node has aggregator")
	assert.Equal(suite.T(), int64(30), trie.root.Children["foo:"].Children["bar:"].Aggr.Params[BytesSize], "Fork node aggregated value is 30")
}

func (suite *TrieTestSuite) TestAddWithMaxChildrenLimit() {
	maxChildrenSavedLimit := 3

	trie := NewTrie(NewPunctuationSplitter(':'), maxChildrenSavedLimit)

	trie.Add("foo:bar:lorem", ParamValue{BytesSize, 10})
	trie.Add("foo:bar:ipsum", ParamValue{BytesSize, 20})
	trie.Add("foo:bar:dolor", ParamValue{BytesSize, 30})
	trie.Add("foo:bar:sit", ParamValue{BytesSize, 40})
	trie.Add("foo:bar:deep:key:nested", ParamValue{BytesSize, 50})
	//trie.Add("foo:bar:", ParamValue{BytesSize, 100})

	// root node
	assert.Equal(suite.T(), int64(150), trie.root.Aggr.Params[BytesSize], "Root node aggregated value is 30")

	// fork node
	assert.NotNil(suite.T(), trie.root.Children["foo:"].Children["bar:"].Aggr, "Fork node has aggregator")
	assert.Len(suite.T(), trie.root.Children["foo:"].Children["bar:"].Children, maxChildrenSavedLimit, "Max children saved is 3")
	assert.Equal(suite.T(), int64(150), trie.root.Children["foo:"].Children["bar:"].Aggr.Params[BytesSize], "Fork node aggregated value is 30")
}

func TestTrieTestSuite(t *testing.T) {
	suite.Run(t, new(TrieTestSuite))
}
