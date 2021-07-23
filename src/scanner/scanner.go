package scanner

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/rs/zerolog"
)

type RedisScanner struct {
	client radix.Client
	logger zerolog.Logger
}

func NewScanner(client radix.Client, logger zerolog.Logger) *RedisScanner {
	return &RedisScanner{
		client: client,
		logger: logger,
	}
}

type ScanOptions struct {
	Pattern          string
	ScanCount        int
	PullRoutineCount int
}

func (s *RedisScanner) Scan(options ScanOptions, result *trie.Trie) {
	var key string
	scanOpts := radix.ScannerConfig{
		Command: "SCAN",
		Count:   options.ScanCount,
	}

	if options.Pattern != "*" {
		scanOpts.Pattern = options.Pattern
	}

	radixScanner := scanOpts.New(s.client)
	for radixScanner.Next(context.Background(), &key) {
		var res int64
		err := s.client.Do(context.Background(), radix.Cmd(&res, "MEMORY", "USAGE", key))
		if err != nil {
			s.logger.Error().Err(err).Msgf("Error dumping key %s", key)
			continue
		}

		result.Add(
			key,
			trie.ParamValue{Param: trie.BytesSize, Value: res},
			trie.ParamValue{Param: trie.KeysCount, Value: 1},
		)
		s.logger.Debug().Msgf("Dump %s value: %d", key, res)
	}
}
