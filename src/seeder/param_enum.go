package seeder

import "math/rand"

type EnumSeedParameter struct {
	values []string
	count  int
}

func NewEnumSeedParameter(values ...string) EnumSeedParameter {
	return EnumSeedParameter{
		values,
		len(values),
	}
}

func (e EnumSeedParameter) generate() string {
	return e.values[rand.Intn(e.count)]
}
