package seeder

type RecordGenerator interface {
	generate() <-chan Record
}

type StringGenerator interface {
	generate() string
}

type IntGenerator interface {
	generate() int
}

type Record struct {
	key   string
	ttl   int
	value string
}

func NewGenericRecordGenerator(repeatCount int, keyGenerator, valueGenerator StringGenerator, ttlGenerator IntGenerator) GenericRecordGenerator {
	return GenericRecordGenerator{
		repeatCount,
		keyGenerator,
		valueGenerator,
		ttlGenerator,
	}
}

type GenericRecordGenerator struct {
	repeatCount    int
	keyGenerator   StringGenerator
	valueGenerator StringGenerator
	ttlGenerator   IntGenerator
}

func (g GenericRecordGenerator) generate() <-chan Record {
	result := make(chan Record)
	go func() {
		defer close(result)
		for i := 0; i < g.repeatCount; i++ {
			result <- Record{
				g.keyGenerator.generate(),
				g.ttlGenerator.generate(),
				g.valueGenerator.generate(),
			}

		}
	}()

	return result
}
