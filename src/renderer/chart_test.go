package renderer

import (
	"encoding/xml"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
	"testing"
)

type ChartRendererTestSuite struct {
	suite.Suite
	trie *trie.Trie
}

func (suite *ChartRendererTestSuite) TestRender() {
	srvMock := &MockServer{}
	srvMock.On("Serve", 123, "test page content").Once()

	rendererMock := &mockPageRenderer{}
	expectedRendererResult := Node{
		Name:       "Total",
		Value:      6500,
		ValueHuman: "6.3K",
		FullPath:   "Total",
		KeysCount:  9,
		Children: []Node{
			{
				Name:       "dev:",
				Value:      2500,
				ValueHuman: "2.4K",
				FullPath:   "dev:",
				KeysCount:  7,
				Children: []Node{
					{
						Name:       "article:",
						Value:      500,
						ValueHuman: "500B",
						FullPath:   "dev:article:",
						KeysCount:  5,
						Children: []Node{
							{Name: "1", Value: 100, ValueHuman: "100B", FullPath: "dev:article:1", KeysCount: 1},
							{Name: "2", Value: 100, ValueHuman: "100B", FullPath: "dev:article:2", KeysCount: 1},
							{Name: "3", Value: 100, ValueHuman: "100B", FullPath: "dev:article:3", KeysCount: 1},
						},
					}, {
						Name:       "user:",
						Value:      2000,
						ValueHuman: "2K",
						FullPath:   "dev:user:",
						KeysCount:  2,
						Children: []Node{
							{Name: "bar", Value: 1000, ValueHuman: "1000B", FullPath: "dev:user:bar", KeysCount: 1},
							{Name: "foo", Value: 1000, ValueHuman: "1000B", FullPath: "dev:user:foo", KeysCount: 1},
						},
					},
				},
			},
			{
				Name:       "prod:user:",
				Value:      4000,
				ValueHuman: "3.9K",
				FullPath:   "prod:user:",
				KeysCount:  2,
				Children: []Node{
					{Name: "bar", Value: 2000, ValueHuman: "2K", FullPath: "prod:user:bar", KeysCount: 1},
					{Name: "foo", Value: 2000, ValueHuman: "2K", FullPath: "prod:user:foo", KeysCount: 1},
				},
			},
		},
	}
	rendererMock.On("render", expectedRendererResult).Once().Return("test page content", nil)

	renderer := NewChartRenderer(
		srvMock,
		ChartRendererParams{
			Depth: 2,
			Port:  123,
		},
	)
	renderer.pageRenderer = rendererMock

	err := renderer.Render(suite.trie.Root())

	suite.Assert().Nil(err)
	srvMock.AssertExpectations(suite.T())
	rendererMock.AssertExpectations(suite.T())
}

func (suite *ChartRendererTestSuite) TestAnychartRenderer() {
	renderer := anychartRenderer{}
	res, err := renderer.render(Node{
		Name:       "Total",
		Value:      6500,
		ValueHuman: "6.3K",
		FullPath:   "Total",
		KeysCount:  9,
		Children: []Node{
			{
				Name:       "dev:",
				Value:      100,
				ValueHuman: "100B",
				FullPath:   "dev:",
				KeysCount:  7,
			},
			{
				Name:       "prod:",
				Value:      2500,
				ValueHuman: "2.4K",
				FullPath:   "dev:",
				KeysCount:  15,
			},
		},
	})

	r := strings.NewReader(res)
	d := xml.NewDecoder(r)

	// Configure the decoder for HTML; leave off strict and autoclose for XHTML
	d.Strict = true
	//d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var errParse error
	for errParse == nil {
		_, errParse = d.Token()
	}

	suite.Assert().Equal(io.EOF, errParse, "Invalid html: %s", errParse)
	suite.Assert().Nil(err)
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

type mockPageRenderer struct {
	mock.Mock
}

func (m *mockPageRenderer) render(result Node) (string, error) {
	args := m.Called(result)

	return args.String(0), args.Error(1)
}
