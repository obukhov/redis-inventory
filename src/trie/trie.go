package trie

import (
	"io"
)

func NewTrie(splitter Splitter, maxBucketSize uint64) *Trie {
	return &Trie{
		root:          NewNode(""),
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
		if !curNode.HasChild(keyPiece) { // can be simplified when working directly with map
			nextNode = NewNode(keyPiece)
			curNode.AddChild(keyPiece, nextNode)
		} else {
			nextNode = curNode.GetChild(keyPiece)
		}

		if curNode.IsFork() && !curNode.HasAggregator() { // no need to check if no new child haven't been added
			curNode.AddAggregator(NewAggregator())
		}

		if curNode.HasAggregator() {
			for _, p := range paramValues {
				curNode.Aggregator().Add(p.Param, p.Value)
			}
		}

		curNode = nextNode
	}
}

func (t *Trie) Dump(w io.Writer) {
	// Write trie summary
}
