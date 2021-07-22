package renderer

import (
	"github.com/hetiansu5/urlquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/obukhov/redis-inventory/src/trie"
	"io"
	"os"
	"strconv"
	"strings"
)

func NewTableRendererParams(paramsString string) (TableRendererParams, error) {
	params := TableRendererParams{}

	err := urlquery.Unmarshal([]byte(paramsString), &params)
	if err != nil {
		return params, err
	}

	params.indent = params.Padding + strings.Repeat(" ", params.PaddingSpaceCount)

	return params, nil
}

type TableRendererParams struct {
	Depth             int    `query:"depth"`
	Padding           string `query:"padding"`
	PaddingSpaceCount int    `query:"padSpaces"`
	indent            string
}

type TableRenderer struct {
	output io.Writer
	params TableRendererParams
}

func (o TableRenderer) Render(trie *trie.Trie) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "ByteSize", "Count"})

	o.appendLevel(t, trie.Root(), 1, "")
	t.Render()

	return nil
}

func (o TableRenderer) appendLevel(t table.Writer, node *trie.Node, level int, prefix string) {
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

		if level < o.params.Depth {
			o.appendLevel(t, childNode, nextLevel, prefix+key)
		}
	}
}

func (o TableRenderer) displayKey(level int, key string, prefix string) string {
	if o.params.indent != "" {
		return strings.Repeat(o.params.indent, level) + key
	} else {
		return prefix + key
	}
}
