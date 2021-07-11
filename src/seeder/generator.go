package seeder

type RecordGenerator interface {
	generate() <-chan Record
}

type Record struct {
	key   string
	ttl   int
	value string
}

func NewGenericRecordGenerator(repeatCount int, keyGenerator, valueGenerator StringGenerator) GenericRecordGenerator {
	return GenericRecordGenerator{
		repeatCount,
		keyGenerator,
		valueGenerator,
	}
}

type GenericRecordGenerator struct {
	repeatCount    int
	keyGenerator   StringGenerator
	valueGenerator StringGenerator
}

func (g GenericRecordGenerator) generate() <-chan Record {
	result := make(chan Record)
	go func() {
		defer close(result)
		for i := 0; i < g.repeatCount; i++ {
			result <- Record{
				key:   g.keyGenerator.generate(),
				ttl:   86400, //todo add generator
				value: g.valueGenerator.generate(),
			}

		}
	}()

	return result
}
