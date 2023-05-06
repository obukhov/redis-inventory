package scanner

import (
	"context"
	"runtime"
	"sync"

	"github.com/obukhov/redis-inventory/src/adapter"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/rs/zerolog"
)

// RedisServiceInterface abstraction to access redis
type RedisServiceInterface interface {
	ScanKeys(ctx context.Context, options adapter.ScanOptions) <-chan string
	GetKeysCount(ctx context.Context) (int64, error)
	GetMemoryUsage(ctx context.Context, key string) (int64, error)
}

// RedisScanner scans redis keys and puts them in a trie
type RedisScanner struct {
	redisService RedisServiceInterface
	scanProgress adapter.ProgressWriter
	logger       zerolog.Logger
}

// NewScanner creates RedisScanner
func NewScanner(redisService RedisServiceInterface, scanProgress adapter.ProgressWriter, logger zerolog.Logger) *RedisScanner {
	return &RedisScanner{
		redisService: redisService,
		scanProgress: scanProgress,
		logger:       logger,
	}
}

// Scan initiates scanning process
func (s *RedisScanner) Scan(options adapter.ScanOptions, result *trie.Trie) {
	var totalCount int64
	if options.Pattern == "*" || options.Pattern == "" {
		totalCount = s.getKeysCount()
	}

	s.scanProgress.Start(totalCount)

	cpus := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(cpus)

	keys := s.redisService.ScanKeys(context.Background(), options)
	for i := 0; i < cpus; i++ {
		go func(batch int) {
			defer wg.Done()

			for key := range keys {
				s.scanProgress.Increment()
				res, err := s.redisService.GetMemoryUsage(context.Background(), key)
				if err != nil {
					s.logger.Error().Err(err).Msgf("Error dumping key %s", key)
					return
				}

				result.Add(
					key,
					trie.ParamValue{Param: trie.BytesSize, Value: res},
					trie.ParamValue{Param: trie.KeysCount, Value: 1},
				)

				s.logger.Debug().Msgf("Dump %s value: %d", key, res)
			}
		}(i)
	}

	wg.Wait()
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
