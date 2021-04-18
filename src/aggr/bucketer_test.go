package aggr

import (
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

type BucketerTestSuite struct {
	suite.Suite
}

func (suite *BucketTestSuite) TestAdd() {
	bucketer := NewBucketer(math.MaxInt64)

	bucketer.Add([]string{"dev", "foo", "1"}, map[string]uint64{"cnt": 1, "size": 100})
	bucketer.Add([]string{"dev", "foo", "2"}, map[string]uint64{"cnt": 1, "size": 10})

	bucketer.Add([]string{"dev", "bar", "3"}, map[string]uint64{"cnt": 1, "size": 50})

	bucketer.Add([]string{"prod", "counter"}, map[string]uint64{"cnt": 1, "size": 1})
	bucketer.Add([]string{"prod", "lock"}, map[string]uint64{"cnt": 1, "size": 2})

	devBucket, devFound := bucketer.Get("dev")

	suite.Require().True(devFound)

	suite.Equal(map[string]uint64{"cnt": 3, "size": 160}, devBucket.GetParams())
	suite.False(devBucket.IsLeaf())

	fooBucket, fooFound := devBucket.Get("foo")

	suite.Require().True(fooFound)

	suite.Equal(map[string]uint64{"cnt": 2, "size": 110}, fooBucket.GetParams())
	suite.False(fooBucket.IsLeaf())

	key1, keyFound := fooBucket.Get("1")
	suite.Require().True(keyFound)

	suite.True(key1.IsLeaf())
	suite.Equal(map[string]uint64{}, key1.GetParams()) // weird

}

func TestBucketerTestSuite(t *testing.T) {
	suite.Run(t, new(BucketerTestSuite))
}
