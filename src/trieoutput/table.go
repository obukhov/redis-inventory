package trieoutput

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/obukhov/redis-inventory/src/trie"
	"io"
	"os"
	"strconv"
	"strings"
)

func NewTableTrieOutput(output io.Writer, depth int, padding string) *TableTrieOutput {
	return &TableTrieOutput{
		output:  output,
		depth:   depth,
		padding: padding,
	}
}

type TableTrieOutput struct {
	output  io.Writer
	depth   int
	padding string
}

func (o *TableTrieOutput) Render(trie *trie.Trie) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "ByteSize", "Count"})

	o.appendLevel(t, trie.Root(), 1, "")
	t.Render()

}

func (o *TableTrieOutput) appendLevel(t table.Writer, node *trie.Node, level int, prefix string) {
	for key, childNode := range node.Children {
		byteSizeColumn := ""
		nextLevel := level + 1
		if !childNode.HasAggregator() {
			var keys []string
			keys, childNode = childNode.FindNextAggregatedNodeWithKey()
			nextLevel += len(keys)
			key = key + strings.Join(keys, "")
		}

		byteSizeColumn = strconv.Itoa(int(childNode.Aggregator().Params[trie.BytesSize]))
		t.AppendRow(table.Row{
			o.displayKey(level, key, prefix),
			byteSizeColumn,
		})

		if level < o.depth {
			o.appendLevel(t, childNode, nextLevel, prefix+key)
		}
	}
}

func (o *TableTrieOutput) displayKey(level int, key string, prefix string) string {
	if o.padding != "" {
		return strings.Repeat(o.padding, level) + key
	} else {
		return prefix + key
	}
}
