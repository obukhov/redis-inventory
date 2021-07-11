package seeder

import "math/rand"

type EnumStringGenerator struct {
	values []string
	count  int
}

func NewEnumStringGenerator(values ...string) EnumStringGenerator {
	return EnumStringGenerator{
		values,
		len(values),
	}
}

func (e EnumStringGenerator) generate() string {
	return e.values[rand.Intn(e.count)]
}
