package seeder

import (
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"log"
)

func NewSeeder(c radix.Client) *Seeder {
	return &Seeder{c}
}

type Seeder struct {
	client radix.Client
}

func (s *Seeder) Seed(patterns ...SeedPattern) {
	for _, p := range patterns {
		for i := 0; i < p.repeatCount; i++ {
			seedParamValues := make([]interface{}, 0, len(p.seedParams))
			for _, seedParam := range p.seedParams {
				seedParamValues = append(seedParamValues, seedParam.generate())
			}
			log.Println(fmt.Sprintf(p.pattern, seedParamValues...)) // todo write to redis
		}
	}
}
