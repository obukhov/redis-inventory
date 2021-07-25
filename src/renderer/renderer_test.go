package renderer

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type RendererTestSuite struct {
	suite.Suite
}

func (suite *RendererTestSuite) TestNewRender() {
	for _, t := range []struct {
		output       string
		outputParams string
		expectedType Renderer
	}{
		{
			"table",
			"depth=2",
			TableRenderer{},
		},
		{
			"json",
			"",
			JSONRenderer{},
		},
	} {
		suite.Run(t.outputParams, func() {
			renderer, err := NewRenderer(t.output, t.outputParams)

			suite.Assert().Nil(err)
			suite.Assert().IsTypef(t.expectedType, renderer, "Unexpected type")
		})
	}
}

func (suite *RendererTestSuite) TestNewRenderError() {
	renderer, err := NewRenderer("foo", "")

	suite.Assert().NotNil(err)
	suite.Assert().Nil(renderer)
}

func TestRendererTestSuite(t *testing.T) {
	suite.Run(t, new(RendererTestSuite))
}
