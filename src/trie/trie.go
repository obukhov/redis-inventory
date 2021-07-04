package trie

import (
	"encoding/json"
	"io"
)

func NewTrie(splitter Splitter, maxBucketSize uint64) *Trie {
	node := NewNode("")
	node.AddAggregator(NewAggregator())

	return &Trie{
		root:          node,
		splitter:      splitter,
		maxBucketSize: maxBucketSize,
	}
}

type Trie struct {
	root          *Node
	splitter      Splitter
	maxBucketSize uint64
}

func (t *Trie) Add(key string, paramValues ...ParamValue) {
	curNode := t.root
	for _, keyPiece := range t.splitter.Split(key) { // change to zero allocation segmenter
		var nextNode *Node

		if childNode := curNode.GetChild(keyPiece); childNode == nil {
			if curNode.ChildCount() == 1 {
				// creating a fork in the trie
				nextAggregatedNode := curNode.FindNextAggregatedNode()
				curNode.AddAggregator(nextAggregatedNode.Aggregator().Clone())
			}

			nextNode = NewNode(keyPiece)
			curNode.AddChild(keyPiece, nextNode)

		} else {
			nextNode = childNode
		}

		if curNode.HasAggregator() {
			for _, p := range paramValues {
				curNode.Aggregator().Add(p.Param, p.Value)
			}
		}

		curNode = nextNode
	}

	curNode.AddAggregator(NewAggregator())
	for _, p := range paramValues {
		curNode.Aggregator().Add(p.Param, p.Value)
	}
}

func (t *Trie) Dump(w io.Writer) {
	e := json.NewEncoder(w)
	e.Encode(t.root)
}
