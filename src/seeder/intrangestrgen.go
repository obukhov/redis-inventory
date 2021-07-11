package seeder

import (
	"math/rand"
	"strconv"
)

type IntRangeStringGenerator struct {
	min int
	max int
}

func NewIntStringGenerator(min, max int) IntRangeStringGenerator {
	return IntRangeStringGenerator{min, max}
}

func (i IntRangeStringGenerator) generate() string {
	return strconv.Itoa(rand.Intn(i.max-i.min) + i.min)
}
