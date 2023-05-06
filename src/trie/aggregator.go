package trie

import "sync"

// ParamValue value for inventory param
type ParamValue struct {
	Param InvParam
	Value int64
}

// NewAggregator creates Aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		Params: make(map[InvParam]int64),
	}
}

// Aggregator struct holding various inventory param values
type Aggregator struct {
	sync.RWMutex
	Params map[InvParam]int64
}

// Add adds inv parameter value to aggregation
func (a *Aggregator) Add(param InvParam, val int64) {
	a.Lock()
	a.Params[param] += val
	a.Unlock()
}

// Clone creates a copy of aggregator
func (a *Aggregator) Clone() *Aggregator {
	a.RLock()
	defer a.RUnlock()

	cloned := NewAggregator()
	for k, v := range a.Params {
		cloned.Params[k] = v
	}

	return cloned
}
