package renderer

import (
	"code.cloudfoundry.org/bytefmt"
	"encoding/json"
	"github.com/hetiansu5/urlquery"
	"github.com/obukhov/redis-inventory/src/trie"
	"io"
	"io/fs"
	"io/ioutil"
	"strings"
)

// NewChartRendererParams creates ChartRendererParams
func NewChartRendererParams(paramsSerialized string) (ChartRendererParams, error) {
	params := ChartRendererParams{Depth: 10}

	err := urlquery.Unmarshal([]byte(paramsSerialized), &params)
	if err != nil {
		return params, err
	}

	return params, nil
}

// ChartRendererParams represents rendering params for web renderer
type ChartRendererParams struct {
	Depth int `query:"depth"`
}

// NewChartRenderer creates ChartRenderer
func NewChartRenderer(output io.Writer, params ChartRendererParams) ChartRenderer {
	return ChartRenderer{
		params: params,
	}
}

// ChartRenderer renders trie in the JSON format
type ChartRenderer struct {
	params ChartRendererParams
}

// Render executes rendering
func (o ChartRenderer) Render(root *trie.Node) error {

	result := o.convertChildren(root, 0, "")
	s, err := json.Marshal(result)
	if err != nil {
		return err
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

				// configure labels
				chart.labels().format("{%name}");

				// configure tooltips
				chart.tooltip().useHtml(true);
				chart.tooltip().format("<span style='font-weight:bold'>{%pathFull}</span><br>{%valueHuman}");

				// initiate drawing the chart
				chart.draw();
			</script>
		</body>
	</html>
	`

	return ioutil.WriteFile("build/chart.html", []byte(rendered), fs.ModePerm)
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

		value := childNode.Aggregator().Params[trie.BytesSize]
		item := Node{
			Name:       key,
			Value:      value,
			ValueHuman: bytefmt.ByteSize(uint64(value)),
			FullPath:   prefix + key,
		}

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

// Node structure for serialized json of anychart library
type Node struct {
	Name       string `json:"name"`
	Value      int64  `json:"value"`
	ValueHuman string `json:"valueHuman"`
	FullPath   string `json:"pathFull"`
	Children   []Node `json:"children"`
}
