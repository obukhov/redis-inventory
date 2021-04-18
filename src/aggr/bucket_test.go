package aggr

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type BucketTestSuite struct {
	suite.Suite
}

func (suite *BucketTestSuite) TestAggregate() {
	for _, testcase := range []struct {
		in      []map[string]uint64
		out     map[string]uint64
		message string
	}{
		{
			[]map[string]uint64{
				{
					"foo": 5,
				},
				{
					"foo": 2,
					"bar": 1,
				},
			},
			map[string]uint64{
				"foo": 7,
				"bar": 1,
			},
			"aggregate missing key",
		},
		{
			[]map[string]uint64{},
			map[string]uint64{},
			"empty map",
		},
	} {
		bucket := NewBucket("")
		for _, v := range testcase.in {
			bucket.Aggregate(v)
		}

		suite.Equal(testcase.out, bucket.GetParams(), testcase.message)
	}
}

func (suite *BucketTestSuite) TestHas() {
	bucket := NewBucket("")

	bucket.Add(NewBucket("foo"))
	bucket.Add(NewBucket("bar"))

	suite.Equal(true, bucket.Has("foo"))
	suite.Equal(false, bucket.Has("zap"))
}

func (suite *BucketTestSuite) TestGet() {
	bucket := NewBucket("")

	foo := NewBucket("foo")
	bucket.Add(foo)
	bucket.Add(NewBucket("bar"))

	getFoo, foundFoo := bucket.Get("foo")

	suite.Equal(foo, getFoo)
	suite.True(foundFoo)

	getZap, foundZap := bucket.Get("zap")

	suite.Nil(getZap)
	suite.False(foundZap)
}

func (suite *BucketTestSuite) TestIsLeaf() {
	bucket := NewBucket("")

	suite.False(bucket.IsLeaf())

	bucket.SetIsLeaf(true)
	suite.True(bucket.IsLeaf())
}

func TestBucketTestSuite(t *testing.T) {
	suite.Run(t, new(BucketTestSuite))
}
