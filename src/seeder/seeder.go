package seeder

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	"log"
)

func NewSeeder(c radix.Client) *Seeder {
	return &Seeder{c}
}

type Seeder struct {
	client radix.Client
}

func (s *Seeder) Seed(generators ...RecordGenerator) {
	for _, generator := range generators {
		for record := range generator.generate() {
			err := s.client.Do(context.Background(), radix.FlatCmd(nil, "SET", record.key, record.value))
			if err != nil {
				log.Printf("Error creating key %s: %s", record.key, err)
				continue
			}

			if record.ttl >= 0 {
				err := s.client.Do(context.Background(), radix.FlatCmd(nil, "EXPIRE", record.key, record.ttl))
				if err != nil {
					log.Printf("Error setting TTL for the key %s: %s", record.key, err)
					continue
				}
			}
		}
	}
}
