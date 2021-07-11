package seeder

import (
	"fmt"
)

func NewPatternStringGenerator(pattern string, paramsGenerator ...StringGenerator) PatternStringGenerator {
	return PatternStringGenerator{
		pattern:         pattern,
		paramsGenerator: paramsGenerator,
	}
}

type PatternStringGenerator struct {
	pattern         string
	paramsGenerator []StringGenerator
}

func (p PatternStringGenerator) generate() string {
	paramValues := make([]interface{}, 0, len(p.paramsGenerator))
	for _, paramGenerator := range p.paramsGenerator {
		paramValues = append(paramValues, paramGenerator.generate())
	}

	return fmt.Sprintf(p.pattern, paramValues...)
}
