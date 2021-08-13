package renderer

import (
	"code.cloudfoundry.org/bytefmt"
	"encoding/json"
	"errors"
	"github.com/hetiansu5/urlquery"
	"github.com/obukhov/redis-inventory/src/server"
	"github.com/obukhov/redis-inventory/src/trie"
	"sort"
	"strings"
)

// NewChartRendererParams creates ChartRendererParams
func NewChartRendererParams(paramsSerialized string) (ChartRendererParams, error) {
	params := ChartRendererParams{Depth: 10, Port: 8888}

	err := urlquery.Unmarshal([]byte(paramsSerialized), &params)
	if err != nil {
		return params, err
	}

	if params.Port <= 0 {
		return params, errors.New("port cannot be negative")
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
		server:       srv,
		pageRenderer: anychartRenderer{},
		params:       params,
	}
}

// ChartRenderer renders trie in the JSON format
type ChartRenderer struct {
	server       server.ServerInterface
	pageRenderer pageRenderer
	params       ChartRendererParams
}

// Render executes rendering
func (o ChartRenderer) Render(root *trie.Node) error {
	result := o.toNode(root, "Total", "")
	result.Children = o.convertChildren(root, 0, "")

	rendered, err := o.pageRenderer.render(result)
	if err != nil {
		return err
	}

	o.server.Serve(o.params.Port, rendered)

	return nil
}

func (o ChartRenderer) convertChildren(node *trie.Node, level int, prefix string) []Node {
	result := make([]Node, 0)

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

		item := o.toNode(childNode, key, prefix)

		if level < o.params.Depth && node.OverflowChildrenCount == 0 {
			item.Children = o.convertChildren(childNode, nextLevel, prefix+key)
		}

		result = append(result, item)
	}

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

type pageRenderer interface {
	render(result Node) (string, error)
}

type anychartRenderer struct{}

func (o anychartRenderer) render(result Node) (string, error) {
	s, err := json.Marshal([]Node{result})
	if err != nil {
		return "", err
	}

	rendered := `<html>
		<head>
			<script src="//unpkg.com/d3"></script>
			<script src="https://cdn.anychart.com/releases/8.10.0/js/anychart-core.min.js"></script>
			<script src="https://cdn.anychart.com/releases/8.10.0/js/anychart-sunburst.min.js"></script>
		</head>
		<body>
			<div id="chart"></div>
			<script type="text/javascript">
				// create data
				var data = ` + string(s) + `;
				const color = d3.scaleOrdinal(d3.schemePaired)
				// Free license provided for the project
				anychart.licenseKey("redis-inventory-80818dc5-535fabc6");
				// create a chart and set the data
				var chart = anychart.sunburst(data, "as-tree");

				// set the container id
				chart.container("chart");
				chart.calculationMode("parent-dependent");

				// configure the visual settings of the chart
				//chart.palette(anychart.palettes.default);
				chart.fill(function () {
				    if (this.level > 1) {
						c = d3.hcl(this.parentColor)
						c.h += (1.5 * (this.index - this.parent.i + 1))
						c.s += 5
						return c.formatHex()
					} else {
					 	return color(this.index);
					}
				});

				chart.labels().position("radial");

				// configure labels
				chart.labels().format("{%name}");

				// configure tooltips
				chart.tooltip().useHtml(true);
				chart.tooltip().format(
					"<span style='font-weight:bold'>{%pathFull}</span><br />{%valueHuman} in {%keys} keys"
				);

				// initiate drawing the chart
				chart.draw();
			</script>
		</body>
	</html>
	`

	return rendered, nil
}
