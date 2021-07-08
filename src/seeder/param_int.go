package seeder

import (
	"math/rand"
	"strconv"
)

type IntSeedParameter struct {
	min int
	max int
}

func NewIntSeedParameter(min, max int) IntSeedParameter {
	return IntSeedParameter{min, max}
}

func (i IntSeedParameter) generate() string {
	return strconv.Itoa(rand.Intn(i.max-i.min) + i.min)
}
