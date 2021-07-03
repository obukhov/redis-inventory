package trie

type InvParam uint

const (
	BytesSize InvParam = iota
)

type ParamValue struct {
	Param InvParam
	Value uint64
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

type Aggregator struct {
	params map[InvParam]uint64
}

func (a *Aggregator) Add(param InvParam, val uint64) {
	a.params[param] += val
}
