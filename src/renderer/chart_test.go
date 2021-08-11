package renderer

import (
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ChartRendererTestSuite struct {
	suite.Suite
	trie *trie.Trie
}

func (suite *ChartRendererTestSuite) TestRender() {
	srvMock := &MockServer{}
	srvMock.On("Serve", 123, mock.Anything).Once()

	renderer := NewChartRenderer(
		srvMock,
		ChartRendererParams{
			Depth: 2,
			Port:  123,
		},
	)

	err := renderer.Render(suite.trie.Root())

	suite.Assert().Nil(err)
	srvMock.AssertExpectations(suite.T())
}

func (suite *ChartRendererTestSuite) SetupTest() {
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

func (suite *ChartRendererTestSuite) setupTrieKey(key string, value int64) {
	suite.trie.Add(key, trie.ParamValue{Param: trie.BytesSize, Value: value}, trie.ParamValue{Param: trie.KeysCount, Value: 1})
}

func TestChartRendererTestSuite(t *testing.T) {
	suite.Run(t, new(ChartRendererTestSuite))
}

type MockServer struct {
	mock.Mock
}

func (m *MockServer) Serve(port int, content string) {
	m.Called(port, content)
}
