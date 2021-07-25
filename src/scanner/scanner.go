package scanner

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/rs/zerolog"
)

// RedisScanner scans redis keys and puts them in a trie
type RedisScanner struct {
	client       radix.Client
	scanProgress ProgressWriter
	logger       zerolog.Logger
}

// NewScanner creates RedisScanner
func NewScanner(client radix.Client, scanProgress ProgressWriter, logger zerolog.Logger) *RedisScanner {
	return &RedisScanner{
		client:       client,
		scanProgress: scanProgress,
		logger:       logger,
	}
}

// ScanOptions options for scanning keyspace
type ScanOptions struct {
	Pattern   string
	ScanCount int
}

// Scan initiates scanning process
func (s *RedisScanner) Scan(options ScanOptions, result *trie.Trie) {

	var key string
	scanOpts := radix.ScannerConfig{
		Command: "SCAN",
		Count:   options.ScanCount,
	}

	var totalCount int64
	if options.Pattern != "*" && options.Pattern != "" {
		scanOpts.Pattern = options.Pattern
	} else {
		totalCount = s.getKeysCount()
	}

	s.scanProgress.Start(totalCount)
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
		s.scanProgress.Increment()
	}
	s.scanProgress.Stop()
}

func (s *RedisScanner) getKeysCount() int64 {
	var keysCount int64
	err := s.client.Do(context.Background(), radix.Cmd(&keysCount, "DBSIZE"))
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Cannot determine number of keys, DBSIZE command returns error")
	}

	return keysCount
}
