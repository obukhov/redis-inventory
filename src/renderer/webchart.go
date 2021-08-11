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

	if rendered, err := o.renderPage(result); err != nil {
		return err
	} else {
		o.server.Serve(o.params.Port, rendered)
		// wait for input
	}

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
	Value      int64  `json:"value"`
	ValueHuman string `json:"valueHuman"`
	FullPath   string `json:"pathFull"`
	KeysCount  int64  `json:"keys"`
	Children   []Node `json:"children"`
}

func (o ChartRenderer) renderPage(result Node) (string, error) {
	s, err := json.Marshal([]Node{result})
	if err != nil {
		return "", err
	}

	rendered := `<html>
		<head>
			<script src="https://cdn.anychart.com/releases/8.10.0/js/anychart-core.min.js"></script>
			<script src="https://cdn.anychart.com/releases/8.10.0/js/anychart-sunburst.min.js"></script>
		</head>
		<body>
			<div id="chart"></div>
			<script type="text/javascript">
				// create data
				var data = ` + string(s) + `;

				// create a chart and set the data
				var chart = anychart.sunburst(data, "as-tree");

				// set the container id
				chart.container("chart");
				chart.calculationMode("parent-dependent");

				// configure the visual settings of the chart
				//chart.palette(anychart.palettes.default);
				chart.fill(function () {
				  return this.parent && this.level > 1 ?
				   anychart.color.lighten(this.parentColor, 0.2) :
				   this.mainColor;
				});

				chart.labels().position("radial");

				// configure labels
				chart.labels().format("{%name}");

				// configure tooltips
				chart.tooltip().useHtml(true);
				chart.tooltip().format(
					"<span style='font-weight:bold'>{%pathFull}</span><br>{%valueHuman} in {%keys} keys"
				);

				// initiate drawing the chart
				chart.draw();
			</script>
		</body>
	</html>
	`

	return rendered, nil
}
