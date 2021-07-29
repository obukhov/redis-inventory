package adapter

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/suite"
	"testing"
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

	client, err := (radix.PoolConfig{}).New(context.Background(), "tcp", m.Addr())

	service := NewRedisService(client)

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

	res := service.ScanKeys(context.Background(), ScanOptions{ScanCount: 1000})

	key1 := <-res
	key2 := <-res

	suite.Assert().Equal("dev:key1", key1)
	suite.Assert().Equal("dev:key2", key2)

	m.Close()

}

func TestRedisServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RedisServiceTestSuite))
}
