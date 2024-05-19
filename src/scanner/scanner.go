package scanner

import (
	"context"

	"github.com/Perseus/redis-inventory/src/adapter"
	"github.com/Perseus/redis-inventory/src/trie"
	"github.com/rs/zerolog"
)

// RedisServiceInterface abstraction to access redis
type RedisServiceInterface interface {
	ScanKeys(ctx context.Context, options adapter.ScanOptions, updateChan chan adapter.ShardScanUpdate)
	ScanAllShards(ctx context.Context, options adapter.ScanOptions, resultChan chan adapter.ShardScanUpdate) chan adapter.ShardInit
	GetKeysCount(ctx context.Context) (int64, error)
	GetMemoryUsage(ctx context.Context, key string) (int64, error)
	GetTotalShards(ctx context.Context) (int64, error)
	GetCurrentNodeIp() string
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
	if options.IsCluster {
		s.ScanCluster(options, result)
	} else {
		s.ScanSingleNode(options, result)
	}
}

func (s *RedisScanner) ScanCluster(options adapter.ScanOptions, result *trie.Trie) {
	scanUpdateChan := make(chan adapter.ShardScanUpdate, 1000)
	totalShards, err := s.redisService.GetTotalShards(context.Background())
	if err != nil {
		s.logger.Error().Err(err).Msg("Error getting total shards")
		return
	}

	shardInitChan := s.redisService.ScanAllShards(context.Background(), options, scanUpdateChan)
	currentShardNum := 0

	for shard := range shardInitChan {
		s.scanProgress.AddTracker(shard.ShardIp, shard.TotalKeys)
		currentShardNum++
		if currentShardNum == int(totalShards) {
			break
		}
	}

	s.scanProgress.SetNumTrackersExpected(int(totalShards))
	s.scanProgress.Start()

	for update := range scanUpdateChan {
		s.scanProgress.Increment(update.ShardIp, 1)
		result.Add(
			update.Key,
			trie.ParamValue{Param: trie.BytesSize, Value: update.MemoryUsage},
			trie.ParamValue{Param: trie.KeysCount, Value: 1},
		)
	}

	s.scanProgress.Stop()
}

func (s *RedisScanner) ScanSingleNode(options adapter.ScanOptions, result *trie.Trie) {
	var totalCount int64
	if options.Pattern == "*" || options.Pattern == "" {
		totalCount = s.getKeysCount()
	}

	nodeIp := s.redisService.GetCurrentNodeIp()
	scanUpdateChan := make(chan adapter.ShardScanUpdate, 1000)

	s.scanProgress.AddTracker(nodeIp, totalCount)
	s.scanProgress.Start()

	s.redisService.ScanKeys(context.Background(), options, scanUpdateChan)

	for update := range scanUpdateChan {
		s.scanProgress.Increment(nodeIp, 1)

		result.Add(
			update.Key,
			trie.ParamValue{Param: trie.BytesSize, Value: update.MemoryUsage},
			trie.ParamValue{Param: trie.KeysCount, Value: 1},
		)

		s.logger.Debug().Msgf("Dump %s value: %d", update.Key, update.MemoryUsage)
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
