package seeder

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/zerolog"
)

func NewSeeder(c radix.Client, logger zerolog.Logger) *Seeder {
	return &Seeder{c, logger}
}

type Seeder struct {
	client radix.Client
	logger zerolog.Logger
}

func (s *Seeder) Seed(generators ...RecordGenerator) {
	keysCount := 0
	for _, generator := range generators {
		for record := range generator.generate() {
			err := s.client.Do(context.Background(), radix.FlatCmd(nil, "SET", record.key, record.value))
			if err != nil {
				s.logger.Error().Err(err).Msgf("Error creating key %s", record.key)
				continue
			}
			keysCount++

			if record.ttl >= 0 {
				err := s.client.Do(context.Background(), radix.FlatCmd(nil, "EXPIRE", record.key, record.ttl))
				if err != nil {
					s.logger.Error().Err(err).Msgf("Error setting TTL for the key %s", record.key)
					continue
				}
			}
		}
	}

	s.logger.Info().Msgf("%d keys set in redis", keysCount)
}
