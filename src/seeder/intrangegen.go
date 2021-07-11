package seeder

import "math/rand"

type IntRangeGenerator struct {
	min int
	max int
}

func NewIntRangeGenerator(min, max int) IntRangeGenerator {
	return IntRangeGenerator{min, max}
}

func (i IntRangeGenerator) generate() int {
	return rand.Intn(i.max-i.min) + i.min
}
