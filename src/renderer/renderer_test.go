package renderer

import (
	"github.com/rs/zerolog"
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
		{
			"chart",
			"",
			ChartRenderer{},
		},
	} {
		suite.Run(t.outputParams, func() {
			renderer, err := NewRenderer(t.output, t.outputParams, zerolog.Nop())

			suite.Assert().Nil(err)
			suite.Assert().IsTypef(t.expectedType, renderer, "Unexpected type")
		})
	}
}

func (suite *RendererTestSuite) TestNewRenderWithError() {
	for _, t := range []struct {
		output       string
		outputParams string
	}{
		{
			"table",
			"padSpaces=asd",
		},
		{
			"json",
			"padSpaces=asd",
		},
		{
			"chart",
			"port=-1",
		},
	} {
		suite.Run(t.outputParams, func() {
			renderer, err := NewRenderer(t.output, t.outputParams, zerolog.Nop())

			suite.Assert().Nil(renderer)
			suite.Assert().Error(err)
		})
	}
}
func (suite *RendererTestSuite) TestNewRenderError() {
	renderer, err := NewRenderer("foo", "", zerolog.Nop())

	suite.Assert().NotNil(err)
	suite.Assert().Nil(renderer)
}

func TestRendererTestSuite(t *testing.T) {
	suite.Run(t, new(RendererTestSuite))
}
