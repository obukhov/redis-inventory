package adapter

import (
	"context"
	"testing"

	"github.com/Perseus/redis-inventory/src/logger"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type RedisServiceTestSuite struct {
	suite.Suite
	service   RedisService
	miniredis *miniredis.Miniredis
}

func (suite *RedisServiceTestSuite) createRedis() (RedisService, *miniredis.Miniredis) {
	var err error
	m, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	m.Set("dev:key1", "bar")
	m.Set("dev:key2", "foobar")

	client := redis.NewClient(&redis.Options{Addr: m.Addr()})
	logger := logger.NewConsoleLogger("trace")
	service := NewRedisService(logger)
	service.SetClient(*client)

	return service, m
}

func (suite *RedisServiceTestSuite) TestCountKeys() {
	service, m := suite.createRedis()

	count, err := service.GetKeysCount(context.Background())

	suite.Assert().Nil(err)
	suite.Assert().Equal(int64(2), count)

	m.Close()
}

func (suite *RedisServiceTestSuite) TestScan() {
	service, m := suite.createRedis()
	scanUpdateChan := make(chan ShardScanUpdate)
	service.ScanKeys(context.Background(), ScanOptions{ScanCount: 1000}, scanUpdateChan)

	key1 := <-scanUpdateChan
	key2 := <-scanUpdateChan

	suite.Assert().Equal("dev:key1", key1.Key)
	suite.Assert().Equal("dev:key2", key2.Key)

	m.Close()
}

func (suite *RedisServiceTestSuite) TestScanMatch() {
	service, m := suite.createRedis()

	scanUpdateChan := make(chan ShardScanUpdate)
	service.ScanKeys(context.Background(), ScanOptions{Pattern: "*:key1", ScanCount: 1000}, scanUpdateChan)

	key1 := <-scanUpdateChan
	suite.Assert().Equal("dev:key1", key1.Key)

	m.Close()
}

func TestRedisServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RedisServiceTestSuite))
}
