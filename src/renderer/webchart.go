package renderer

import (
	"code.cloudfoundry.org/bytefmt"
	"encoding/json"
	"github.com/hetiansu5/urlquery"
	"github.com/obukhov/redis-inventory/src/server"
	"github.com/obukhov/redis-inventory/src/trie"
	"strings"
)

// NewChartRendererParams creates ChartRendererParams
func NewChartRendererParams(paramsSerialized string) (ChartRendererParams, error) {
	params := ChartRendererParams{Depth: 10, Port: 8888}

	err := urlquery.Unmarshal([]byte(paramsSerialized), &params)
	if err != nil {
		return params, err
	}

	return params, nil
}

// ChartRendererParams represents rendering params for web renderer
type ChartRendererParams struct {
	Depth int `query:"depth"`
	Port  int `query:"port"`
}

// NewChartRenderer creates ChartRenderer
func NewChartRenderer(srv server.ServerInterface, params ChartRendererParams) ChartRenderer {
	return ChartRenderer{
		server: srv,
		params: params,
	}
}

// ChartRenderer renders trie in the JSON format
type ChartRenderer struct {
	server server.ServerInterface
	params ChartRendererParams
}

// Render executes rendering
func (o ChartRenderer) Render(root *trie.Node) error {
	result := o.toNode(root, "Total", "")
	result.Children = o.convertChildren(root, 0, "")

	rendered, err := o.renderPage(result)
	if err != nil {
		return err
	}

	o.server.Serve(o.params.Port, rendered)

	return nil
}

func (o ChartRenderer) convertChildren(node *trie.Node, level int, prefix string) []Node {
	result := make([]Node, 0)

	for key, childNode := range node.Children {

		nextLevel := level + 1
		if !childNode.HasAggregator() {
			var keys []string
			keys, childNode = childNode.FindNextAggregatedNodeWithKey()
			nextLevel += len(keys)
			key = key + strings.Join(keys, "")
		}

		item := o.toNode(childNode, key, prefix)

		if level < o.params.Depth && node.OverflowChildrenCount == 0 {
			item.Children = o.convertChildren(childNode, nextLevel, prefix+key)
		}

		result = append(result, item)
	}

	//if node.OverflowChildrenCount > 0 {
	//	t.AppendRow(table.Row{o.displayKey(level, fmt.Sprintf("( %d more keys )", node.OverflowChildrenCount), prefix)})
	//}

	return result
}

func (o ChartRenderer) toNode(childNode *trie.Node, key string, prefix string) Node {
	//var value int64 = 0
	//if len(childNode.Children) == 0 || childNode.OverflowChildrenCount > 0 {
	//}

	value := childNode.Aggregator().Params[trie.BytesSize]
	item := Node{
		Name:       key,
		Value:      value,
		KeysCount:  childNode.Aggregator().Params[trie.KeysCount],
		ValueHuman: bytefmt.ByteSize(uint64(value)),
		FullPath:   prefix + key,
	}
	return item
}

// Node structure for serialized json of anychart library
type Node struct {
	Name       string `json:"name"`
	Value      int64  `json:"value,omitempty"`
	ValueHuman string `json:"valueHuman"`
	FullPath   string `json:"pathFull"`
	KeysCount  int64  `json:"keys"`
	Children   []Node `json:"children"`
}

func (o ChartRenderer) renderPage(result Node) (string, error) {
	//s, err := json.MarshalIndent(result, "", "  ")
	s, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	rendered := `<html>
		<head>
			<script src="//unpkg.com/d3"></script>
			<script src="//unpkg.com/sunburst-chart"></script>
		</head>
		<body>
			<div id="chart"></div>
			<script type="text/javascript">
				const data = ` + string(s) + `;
 				const color = d3.scaleOrdinal(d3.schemePaired);
				const myChart = Sunburst();
				myChart
					.data(data)
					.label('name')
					.size('value')
					.excludeRoot(false)
					.centerRadius(0)
					.radiusScaleExponent(1)
					.labelOrientation('radial')
				    .color((d, parent) => {
						if (parent && parent.depth > 0) {
							var c = d3.hsl(parent.data.color)
							c.h += 20
							d.color = c + ""
						} else {
							d.color = color(d.name)
						}

						return d.color
					})
					.tooltipContent((d, node) => ` + "`Size: <i>${d.valueHuman}</i> in ${d.keys} key(s)`" + `)
					(document.getElementById('chart'));
			</script>
		</body>
	</html>
	`

	return rendered, nil
}
