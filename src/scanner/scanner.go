package scanner

import (
	"context"

	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/rs/zerolog"
)

// RedisScanner scans redis keys and puts them in a trie
type RedisScanner struct {
	redisService RedisServiceInterface
	scanProgress ProgressWriter
	logger       zerolog.Logger
}

// NewScanner creates RedisScanner
func NewScanner(redisService RedisServiceInterface, scanProgress ProgressWriter, logger zerolog.Logger) *RedisScanner {
	return &RedisScanner{
		redisService: redisService,
		scanProgress: scanProgress,
		logger:       logger,
	}
}

// Scan initiates scanning process
func (s *RedisScanner) Scan(options ScanOptions, result *trie.Trie) {
	var totalCount int64
	if options.Pattern == "*" || options.Pattern == "" {
		totalCount = s.getKeysCount()
	}

	s.scanProgress.Start(totalCount)
	for key := range s.redisService.ScanKeys(context.Background(), options) {
		res, err := s.redisService.GetMemoryUsage(context.Background(), key)
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
	res, err := s.redisService.GetKeysCount(context.Background())
	if err != nil {
		s.logger.Error().Err(err).Msgf("Error getting number of keys")
		return 0
	}

	return res
}
