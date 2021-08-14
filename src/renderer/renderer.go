package renderer

import (
	"errors"
	"github.com/obukhov/redis-inventory/src/server"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/rs/zerolog"
	"os"
)

// Renderer abstraction for rendering trie to a given output
type Renderer interface {
	// Render executes rendering
	Render(root *trie.Node) error
}

// NewRenderer creates Renderer implementation by type and set of params
func NewRenderer(output, paramsString string, logger zerolog.Logger) (Renderer, error) {
	switch output {
	case "table":
		params, err := NewTableRendererParams(paramsString)
		if err != nil {
			return nil, err
		}

		return TableRenderer{os.Stdout, params}, nil
	case "json":
		params, err := NewJSONRendererParams(paramsString)
		if err != nil {
			return nil, err
		}

		return JSONRenderer{os.Stdout, params}, nil

	case "chart":
		params, err := NewChartRendererParams(paramsString)
		if err != nil {
			return nil, err
		}

		return NewChartRenderer(server.NewServer(logger), params), nil
	default:
		return nil, errors.New("unknown render type")
	}
}
