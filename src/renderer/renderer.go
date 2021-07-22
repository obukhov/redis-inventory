package renderer

import (
	"github.com/obukhov/redis-inventory/src/trie"
	"os"
)

type Renderer interface {
	Render(trie *trie.Trie) error
}

func NewRenderer(output, paramsString string) (Renderer, error) {
	switch output {
	case "table":
		params, err := NewTableRendererParams(paramsString)
		if err != nil {
			return nil, err
		}

		return TableRenderer{os.Stdout, params}, nil
	case "json":
		params, err := NewJsonRendererParams(paramsString)
		if err != nil {
			return nil, err
		}

		return JsonRenderer{os.Stdout, params}, nil
	default:
		panic("Unknown output format: " + output)
	}
}
