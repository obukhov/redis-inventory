package trieoutput

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/obukhov/redis-inventory/src/trie"
	"io"
	"os"
	"strconv"
	"strings"
)

func NewTableTrieOutput(output io.Writer, depth int) *TableTrieOutput {
	return &TableTrieOutput{
		output: output,
		depth:  depth,
	}
}

type TableTrieOutput struct {
	output io.Writer
	depth  int
}

func (o *TableTrieOutput) Render(trie *trie.Trie) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "ByteSize", "Count"})

	o.appendLevel(t, trie.Root(), 1)
	t.Render()

}

func (o *TableTrieOutput) appendLevel(t table.Writer, node *trie.Node, level int) {

	for key, childNode := range node.Children {
		byteSizeColumn := ""
		if childNode.HasAggregator() {
			byteSizeColumn = strconv.Itoa(int(childNode.Aggregator().Params[trie.BytesSize]))
		}

		t.AppendRow(table.Row{
			strings.Repeat("  ", level) + key,
			byteSizeColumn,
		})

		if level < o.depth {
			o.appendLevel(t, childNode, level+1)
		}
	}
}
