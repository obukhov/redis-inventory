package renderer

import (
	"encoding/json"
	"github.com/hetiansu5/urlquery"
	"github.com/obukhov/redis-inventory/src/trie"
	"io"
	"strings"
)

func NewJsonRendererParams(paramsSerialized string) (JsonRendererParams, error) {
	params := JsonRendererParams{}

	err := urlquery.Unmarshal([]byte(paramsSerialized), &params)
	if err != nil {
		return params, err
	}

	return params, nil
}

type JsonRendererParams struct {
	Padding           string `query:"padding"`
	PaddingSpaceCount int    `query:"padSpaces"`
}

type JsonRenderer struct {
	output io.Writer
	params JsonRendererParams
}

func (o JsonRenderer) Render(trie *trie.Trie) error {
	encoder := json.NewEncoder(o.output)

	indent := o.params.Padding + strings.Repeat(" ", o.params.PaddingSpaceCount)
	if indent != "" {
		encoder.SetIndent("", indent)
	}

	return encoder.Encode(trie.Root())
}
