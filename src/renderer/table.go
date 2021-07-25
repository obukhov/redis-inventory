package renderer

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/hetiansu5/urlquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/obukhov/redis-inventory/src/trie"
	"golang.org/x/text/message"
	"io"
	"sort"
	"strconv"
	"strings"
)

var p = message.NewPrinter(message.MatchLanguage("en"))

// NewTableRendererParams Creates parameters structure from url-encoded string respecting some defaults
func NewTableRendererParams(paramsString string) (TableRendererParams, error) {
	params := TableRendererParams{
		Depth: 10,
	}

	err := urlquery.Unmarshal([]byte(paramsString), &params)
	if err != nil {
		return params, err
	}

	params.indent = params.Padding + strings.Repeat(" ", params.PaddingSpaceCount)

	return params, nil
}

// TableRendererParams represents renderer parameters
type TableRendererParams struct {
	Depth             int    `query:"depth"`
	Padding           string `query:"padding"`
	PaddingSpaceCount int    `query:"padSpaces"`
	HumanReadable     bool   `query:"human"`
	indent            string
}

// TableRenderer renders trie as ascii-table to output (most probably stdout)
type TableRenderer struct {
	output io.Writer
	params TableRendererParams
}

// Render executes rendering
func (o TableRenderer) Render(trie *trie.Trie) error {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Key", "ByteSize", "KeysCount"})
	o.appendLevel(t, trie.Root(), 1, "")

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1},
		{Number: 2, Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Number: 3, Align: text.AlignRight, AlignHeader: text.AlignCenter},
	})
	t.SetOutputMirror(o.output)
	t.Render()

	return nil
}

func (o TableRenderer) appendLevel(t table.Writer, node *trie.Node, level int, prefix string) {
	childKeys := make([]string, 0, len(node.Children))
	for k := range node.Children {
		childKeys = append(childKeys, k)
	}

	sort.Strings(childKeys)
	for _, key := range childKeys {
		childNode := node.Children[key]
		nextLevel := level + 1
		if !childNode.HasAggregator() {
			var keys []string
			keys, childNode = childNode.FindNextAggregatedNodeWithKey()
			nextLevel += len(keys)
			key = key + strings.Join(keys, "")
		}

		t.AppendRow(table.Row{
			o.displayKey(level, key, prefix),
			o.formatBytes(childNode.Aggregator().Params[trie.BytesSize]),
			o.formatNumber(childNode.Aggregator().Params[trie.KeysCount]),
		})

		if level < o.params.Depth {
			o.appendLevel(t, childNode, nextLevel, prefix+key)
		}
	}

	if node.OverflowChildrenCount > 0 {
		t.AppendRow(table.Row{o.displayKey(level, fmt.Sprintf("( %d more keys )", node.OverflowChildrenCount), prefix)})
	}
}

func (o TableRenderer) formatBytes(value int64) string {
	if o.params.HumanReadable {
		return bytefmt.ByteSize(uint64(value))
	}

	return strconv.Itoa(int(value))
}

func (o TableRenderer) formatNumber(value int64) string {

	if o.params.HumanReadable {
		return p.Sprint(value)
	}

	return strconv.Itoa(int(value))
}

func (o TableRenderer) displayKey(level int, key string, prefix string) string {
	if o.params.indent != "" {
		return strings.Repeat(o.params.indent, level) + key
	}

	return prefix + key
}
