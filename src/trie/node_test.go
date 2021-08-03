package trie

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type NodeTestSuite struct {
	suite.Suite
}

func (suite *NodeTestSuite) TestFirstChildPanic() {
	node := NewNode()
	suite.Assert().Panics(func() {
		node.FirstChild()
	})
}

func (suite *NodeTestSuite) TestFirstChildWithKeyPanic() {
	node := NewNode()
	suite.Assert().Panics(func() {
		node.FirstChildWithKey()
	})
}

func (suite *NodeTestSuite) TestFindNextAggregatedNode() {
	node := NewNode()

	level1 := NewNode()
	level2 := NewNode()
	level3one := NewNode()
	level3two := NewNode()

	level2.AddChild("one", level3one)
	level2.AddChild("two", level3two)
	level2.AddAggregator(NewAggregator())

	level1.AddChild("bar", level2)
	node.AddChild("foo", level1)

	actual := node.FindNextAggregatedNode()

	suite.Assert().Equal(level2, actual)
}

func (suite *NodeTestSuite) TestFindNextAggregatedNodeWithPath() {
	node := NewNode()

	level1 := NewNode()
	level2 := NewNode()
	level3one := NewNode()
	level3two := NewNode()

	level2.AddChild("one", level3one)
	level2.AddChild("two", level3two)
	level2.AddAggregator(NewAggregator())

	level1.AddChild("bar", level2)
	node.AddChild("foo", level1)

	actгalKeys, actual := node.FindNextAggregatedNodeWithKey()

	suite.Assert().Equal(level2, actual)
	suite.Assert().Equal([]string{"foo", "bar"}, actгalKeys)
}
func TestNodeTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}
