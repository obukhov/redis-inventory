package trie

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PunctuationSplitterTestSuite struct {
	suite.Suite
	splitter Splitter
}

func (suite *PunctuationSplitterTestSuite) SetupTest() {
	suite.splitter = NewPunctuationSplitter([]rune{'_', ':'})
}

func (suite *PunctuationSplitterTestSuite) TestSplit() {
	for _, testcase := range []struct {
		in      string
		out     []string
		message string
	}{
		{
			"helloWorld",
			[]string{"helloWorld"},
			"no punctuation",
		},
		{
			"hello_world",
			[]string{"hello_", "world"},
			"basic split",
		},
		{
			"hello___world",
			[]string{"hello___", "world"},
			"multiple punctuation in a single split",
		},
		{
			"hello_my:dear_:world",
			[]string{"hello_", "my:", "dear_:", "world"},
			"basic split",
		},
		{
			"_hello_world",
			[]string{"_", "hello_", "world"},
			"starting with punctuation",
		},
		{
			"___hello_world",
			[]string{"___", "hello_", "world"},
			"starting with multiple punctuation",
		},
		{
			"hello_world___",
			[]string{"hello_", "world___"},
			"ending with multiple punctuation",
		},
		{
			"_",
			[]string{"_"},
			"just punctuation",
		},
		{
			"",
			[]string{""},
			"empty string",
		},
	} {
		suite.Equal(suite.splitter.Split(testcase.in), testcase.out, testcase.message)
	}
}

func TestPunctuationSplitterTestSuite(t *testing.T) {
	suite.Run(t, new(PunctuationSplitterTestSuite))
}
