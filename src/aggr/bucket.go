package aggr

func NewBucket(key string) *Bucket {
	return &Bucket{
		key:        key,
		subBuckets: make([]*Bucket, 0),
		keyToIndex: make(map[string]uint),
		params:     make(map[string]uint64),
		isLeaf:     false,
	}
}

type Bucket struct {
	key        string
	subBuckets []*Bucket
	keyToIndex map[string]uint
	params     map[string]uint64
	isLeaf     bool
}

func (b *Bucket) Has(node string) bool {
	_, found := b.keyToIndex[node]

	return found
}

func (b *Bucket) Aggregate(params map[string]uint64) {
	for k, v := range params {
		b.params[k] += v
	}
}

func (b *Bucket) GetParams() map[string]uint64 {
	return b.params
}

func (b *Bucket) SetIsLeaf(isLeaf bool) {
	b.isLeaf = isLeaf
}

func (b *Bucket) IsLeaf() bool {
	return b.isLeaf
}

func (b *Bucket) Add(next *Bucket) {
	b.subBuckets = append(b.subBuckets, next)
	b.keyToIndex[next.Key()] = uint(len(b.subBuckets) - 1)
}

func (b *Bucket) Get(node string) (*Bucket, bool) {
	bucketIndex, found := b.keyToIndex[node]
	if !found {
		return nil, false
	}

	return b.subBuckets[bucketIndex], true
}

func (b *Bucket) Key() string {
	return b.key
}
