package trie

import "strconv"

type InvParam uint

const (
	BytesSize InvParam = iota
)

func (p InvParam) String() string {
	switch p {
	case BytesSize:
		return "BytesSize"
	}

	panic("Unknown InvParam: " + strconv.Itoa(int(p)))
}

func (p InvParam) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

type ParamValue struct {
	Param InvParam
	Value int64
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		Params: make(map[InvParam]int64),
	}
}

type Aggregator struct {
	Params map[InvParam]int64
}

func (a *Aggregator) Add(param InvParam, val int64) {
	a.Params[param] += val
}

func (a *Aggregator) Clone() *Aggregator {
	cloned := NewAggregator()
	for k, v := range a.Params {
		cloned.Params[k] = v
	}

	return cloned
}
