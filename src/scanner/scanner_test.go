package scanner

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net"
	"testing"
)

type ScannerTestSuite struct {
	suite.Suite
}

type RedisClientMock struct {
}

func (r *RedisClientMock) Addr() net.Addr {
	return nil
}

func (r *RedisClientMock) Do(context.Context, radix.Action) error {
	return nil
}

func (r *RedisClientMock) Close() error {

	return nil
}

func (suite *ScannerTestSuite) TestScan() {
	_ = NewScanner(&RedisClientMock{}, NewPrettyProgressWriter(ioutil.Discard), zerolog.Nop())
}

func TestScannerTestSuite(t *testing.T) {
	suite.Run(t, new(ScannerTestSuite))
}
