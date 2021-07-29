package scanner

import (
	"context"
	"errors"
	"github.com/obukhov/redis-inventory/src/adapter"
	"testing"

	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ScannerTestSuite struct {
	suite.Suite
}

type RedisServiceMock struct {
	mock.Mock
}

func (m *RedisServiceMock) ScanKeys(ctx context.Context, options adapter.ScanOptions) <-chan string {
	args := m.Called(ctx, options)
	return args.Get(0).(chan string)
}

func (m *RedisServiceMock) GetKeysCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}

func (m *RedisServiceMock) GetMemoryUsage(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return int64(args.Int(0)), args.Error(1)
}

type ProgressWriterMock struct {
	mock.Mock
}

func (m *ProgressWriterMock) Start(total int64) {
	m.Called(total)
}

func (m *ProgressWriterMock) Increment() {
	m.Called()
}

func (m *ProgressWriterMock) Stop() {
	m.Called()
}

func (suite *ScannerTestSuite) TestScan() {
	scanChannel := make(chan string, 5)

	redisMock := &RedisServiceMock{}
	redisMock.
		On("GetKeysCount", mock.Anything).Return(2, nil).
		On(
			"ScanKeys",
			mock.Anything,
			adapter.ScanOptions{
				ScanCount: 1000,
				Throttle:  0,
			},
		).
		Return(scanChannel).
		On("GetMemoryUsage", mock.Anything, "key1").Return(1, nil).
		On("GetMemoryUsage", mock.Anything, "key2").Return(10, nil)

	progressMock := &ProgressWriterMock{}
	progressMock.
		On("Start", int64(2)).Once().
		On("Stop").Once().
		On("Increment").Times(2)

	scanChannel <- "key1"
	scanChannel <- "key2"
	close(scanChannel)

	scanner := NewScanner(redisMock, progressMock, zerolog.Nop())
	scanner.Scan(
		adapter.ScanOptions{
			ScanCount: 1000,
			Throttle:  0,
		},
		trie.NewTrie(trie.NewPunctuationSplitter(':'), 5),
	)

	redisMock.AssertExpectations(suite.T())
	progressMock.AssertExpectations(suite.T())
}

func (suite *ScannerTestSuite) TestScanWithPattern() {
	scanChannel := make(chan string, 5)

	redisMock := &RedisServiceMock{}
	redisMock.
		On(
			"ScanKeys",
			mock.Anything,
			adapter.ScanOptions{
				ScanCount: 1000,
				Throttle:  0,
				Pattern:   "dev:*",
			},
		).
		Return(scanChannel).
		On("GetMemoryUsage", mock.Anything, "key1").Return(1, nil).
		On("GetMemoryUsage", mock.Anything, "key2").Return(10, nil)

	progressMock := &ProgressWriterMock{}
	progressMock.
		On("Start", int64(0)).Once().
		On("Stop").Once().
		On("Increment").Times(2)

	scanChannel <- "key1"
	scanChannel <- "key2"
	close(scanChannel)

	scanner := NewScanner(redisMock, progressMock, zerolog.Nop())
	scanner.Scan(
		adapter.ScanOptions{
			ScanCount: 1000,
			Throttle:  0,
			Pattern:   "dev:*",
		},
		trie.NewTrie(trie.NewPunctuationSplitter(':'), 5),
	)

	redisMock.AssertExpectations(suite.T())
	progressMock.AssertExpectations(suite.T())
}

func (suite *ScannerTestSuite) TestScanWithError() {
	scanChannel := make(chan string, 5)

	redisMock := &RedisServiceMock{}
	redisMock.
		On("GetKeysCount", mock.Anything).Return(2, nil).
		On("ScanKeys", mock.Anything, mock.Anything).Return(scanChannel).
		On("GetMemoryUsage", mock.Anything, "key1").Return(1, errors.New("cannot get memory")).
		On("GetMemoryUsage", mock.Anything, "key2").Return(10, nil)

	progressMock := &ProgressWriterMock{}
	progressMock.
		On("Start", int64(2)).Once().
		On("Stop").Once().
		On("Increment").Times(2)

	scanChannel <- "key1"
	scanChannel <- "key2"
	close(scanChannel)

	scanner := NewScanner(redisMock, progressMock, zerolog.Nop())
	result := trie.NewTrie(trie.NewPunctuationSplitter(':'), 5)
	scanner.Scan(
		adapter.ScanOptions{
			ScanCount: 1000,
			Throttle:  0,
		},
		result,
	)

	redisMock.AssertExpectations(suite.T())
	progressMock.AssertExpectations(suite.T())

	suite.Assert().Equal(int64(10), result.Root().Aggregator().Params[trie.BytesSize])
}

func (suite *ScannerTestSuite) TestScanCantGetCountKeys() {
	scanChannel := make(chan string, 5)

	redisMock := &RedisServiceMock{}
	redisMock.
		On("GetKeysCount", mock.Anything).Return(2, errors.New("cannot get count keys")).
		On("ScanKeys", mock.Anything, mock.Anything).Return(scanChannel).
		On("GetMemoryUsage", mock.Anything, "key1").Return(1, nil).
		On("GetMemoryUsage", mock.Anything, "key2").Return(10, nil)

	progressMock := &ProgressWriterMock{}
	progressMock.
		On("Start", int64(0)).Once().
		On("Stop").Once().
		On("Increment").Times(2)

	scanChannel <- "key1"
	scanChannel <- "key2"
	close(scanChannel)

	scanner := NewScanner(redisMock, progressMock, zerolog.Nop())
	scanner.Scan(
		adapter.ScanOptions{
			ScanCount: 1000,
			Throttle:  0,
		},
		trie.NewTrie(trie.NewPunctuationSplitter(':'), 5),
	)

	redisMock.AssertExpectations(suite.T())
	progressMock.AssertExpectations(suite.T())
}

func TestScannerTestSuite(t *testing.T) {
	suite.Run(t, new(ScannerTestSuite))
}
