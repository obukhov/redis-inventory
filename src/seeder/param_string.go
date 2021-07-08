package seeder

import "math/rand"

type StringSeedParameter struct {
	minLength    int
	maxLength    int
	symbols      []rune
	symbolsCount int
}

func NewStringSeedParameter(minLength, maxLength int, symbols ...rune) StringSeedParameter {
	return StringSeedParameter{
		minLength,
		maxLength,
		symbols,
		len(symbols),
	}
}

func (s StringSeedParameter) generate() string {
	length := rand.Intn(s.maxLength-s.minLength) + s.minLength
	result := make([]rune, 0, length)
	for i := 0; i < length; i++ {
		result = append(result, s.symbols[rand.Intn(s.symbolsCount)])
	}

	return string(result)
}
