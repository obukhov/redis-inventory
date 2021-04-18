package aggr

import "math"

func NewBucketer(maxBucketSize uint64) *Bucketer {
	return &Bucketer{
		root:          NewBucket(""),
		maxBucketSize: math.MaxInt64,
	}
}

type Bucketer struct {
	maxBucketSize uint
	root          *Bucket
}

func (b *Bucketer) Add(path []string, params map[string]uint64) {
	maxDepth := len(path) - 1
	cur := b.root
	for n, node := range path {
		cur.Aggregate(params)
		if next, found := cur.Get(node); found {
			cur = next
		} else {
			newBucket := NewBucket(node)
			cur.Add(newBucket)
			cur = newBucket
		}

		if n == maxDepth {
			cur.SetIsLeaf(true)
		}
	}
}

func (b *Bucketer) Get(node string) (*Bucket, bool) {
	return b.root.Get(node)
}
