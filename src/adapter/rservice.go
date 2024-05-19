package adapter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// ScanOptions options for scanning keyspace
type ScanOptions struct {
	Pattern   string
	ScanCount int
	Throttle  int
	IsCluster bool
}

type ShardScanUpdate struct {
	ShardIp     string
	Key         string
	MemoryUsage int64
}

type ShardInit struct {
	ShardIp   string
	TotalKeys int64
}

// RedisService implementation for iteration over redis
type RedisService struct {
	client        redis.Client
	clusterClient redis.ClusterClient
	isCluster     bool
	logger        zerolog.Logger
}

func NewRedisService(logger zerolog.Logger) RedisService {
	rs := RedisService{
		logger: logger,
	}
	return rs
}

func (rs *RedisService) SetClient(client redis.Client) {
	rs.client = client
}

func (rs *RedisService) SetClusterClient(client redis.ClusterClient) {
	rs.clusterClient = client
	rs.isCluster = true
}

// ScanKeys scans keys asynchroniously and sends them to the returned channel
func (s RedisService) ScanKeys(ctx context.Context, options ScanOptions, resultChan chan ShardScanUpdate) {
	go func() {
		scanIter := s.client.Scan(ctx, 0, options.Pattern, int64(options.ScanCount)).Iterator()
		for scanIter.Next(ctx) {
			key := scanIter.Val()
			memoryUsage, err := s.GetMemoryUsage(ctx, key)
			if err != nil {
				s.logger.Error().Err(err).Msg("Error getting memory usage")
			}
			resultChan <- ShardScanUpdate{ShardIp: s.client.Options().Addr, Key: key, MemoryUsage: memoryUsage}
			if options.Throttle > 0 {
				time.Sleep(time.Nanosecond * time.Duration(options.Throttle))
			}
		}
	}()
}

func (s RedisService) ScanAllShards(ctx context.Context, options ScanOptions, resultChan chan ShardScanUpdate) chan ShardInit {
	shardListChan := make(chan ShardInit)

	go func() {
		s.clusterClient.ForEachMaster(ctx, func(ctx context.Context, master *redis.Client) error {
			totalKeys, err := master.DBSize(ctx).Result()
			if err != nil {
				panic(err)
			}
			shardListChan <- ShardInit{ShardIp: master.Options().Addr, TotalKeys: totalKeys}

			scanIter := master.Scan(ctx, 0, options.Pattern, int64(options.ScanCount)).Iterator()
			for scanIter.Next(ctx) {
				key := scanIter.Val()
				memoryUsage, err := s.GetMemoryUsage(ctx, key)
				if err != nil {
					panic(err)
				}

				resultChan <- ShardScanUpdate{ShardIp: master.Options().Addr, Key: key, MemoryUsage: memoryUsage}
				if options.Throttle > 0 {
					time.Sleep(time.Nanosecond * time.Duration(options.Throttle))
				}
			}

			return nil
		})
	}()

	return shardListChan
}

// GetKeysCount returns number of keys in the current database
func (s RedisService) GetKeysCount(ctx context.Context) (int64, error) {
	if s.isCluster {
		size, err := s.clusterClient.DBSize(ctx).Result()
		if err != nil {
			return 0, err
		}
		return size, nil
	}

	size, err := s.client.DBSize(ctx).Result()
	if err != nil {
		return 0, err
	}

	return size, nil
}

// GetMemoryUsage returns memory usage of given key
func (s RedisService) GetMemoryUsage(ctx context.Context, key string) (int64, error) {
	if s.isCluster {
		size, err := s.clusterClient.MemoryUsage(ctx, key).Result()
		if err != nil {
			return 0, err
		}
		return size, nil
	}

	size, err := s.client.MemoryUsage(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return size, nil
}

func (s RedisService) GetTotalShards(ctx context.Context) (int64, error) {
	if s.isCluster {
		nodes, err := s.clusterClient.ClusterShards(ctx).Result()
		if err != nil {
			return 0, err
		}
		return int64(len(nodes)), nil
	}

	return 1, nil
}

func (s RedisService) GetCurrentNodeIp() string {
	return s.client.Options().Addr
}
