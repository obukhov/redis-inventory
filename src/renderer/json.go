package renderer

import (
	"encoding/json"
	"github.com/hetiansu5/urlquery"
	"github.com/obukhov/redis-inventory/src/trie"
	"io"
	"strings"
)

// NewJSONRendererParams creates JSONRendererParams
func NewJSONRendererParams(paramsSerialized string) (JSONRendererParams, error) {
	params := JSONRendererParams{}

	err := urlquery.Unmarshal([]byte(paramsSerialized), &params)
	if err != nil {
		return params, err
	}

	return params, nil
}

// JSONRendererParams represents rendering params fr Json renderer
type JSONRendererParams struct {
	Padding           string `query:"padding"`
	PaddingSpaceCount int    `query:"padSpaces"`
}

// NewJSONRenderer creates JSONRenderer
func NewJSONRenderer(output io.Writer, params JSONRendererParams) JSONRenderer {
	return JSONRenderer{
		output: output,
		params: params,
	}
}

// JSONRenderer renders trie in the JSON format
type JSONRenderer struct {
	output io.Writer
	params JSONRendererParams
}

// Render executes rendering
func (o JSONRenderer) Render(root *trie.Node) error {
	encoder := json.NewEncoder(o.output)

	indent := o.params.Padding + strings.Repeat(" ", o.params.PaddingSpaceCount)
	if indent != "" {
		encoder.SetIndent("", indent)
	}

	return encoder.Encode(root)
}
