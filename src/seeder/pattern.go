package seeder

type SeedParameter interface {
	generate() string
}

func NewSeedPattern(repeatCount int, pattern string, seedParams ...SeedParameter) SeedPattern {
	return SeedPattern{repeatCount, pattern, seedParams}
}

type SeedPattern struct {
	repeatCount int
	pattern     string
	seedParams  []SeedParameter
}
