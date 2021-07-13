package trie

func NewTrie(splitter Splitter, maxChildren int) *Trie {
	node := NewNode()
	node.AddAggregator(NewAggregator())

	return &Trie{
		root:        node,
		splitter:    splitter,
		maxChildren: maxChildren,
	}
}

type Trie struct {
	root        *Node
	splitter    Splitter
	maxChildren int
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

			if curNode.ChildCount() < t.maxChildren {
				nextNode = NewNode()
				curNode.AddChild(keyPiece, nextNode)
			} else {
				break
			}

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

	if !curNode.HasAggregator() {
		if curNode.HasChildren() {
			curNode.AddAggregator(curNode.FindNextAggregatedNode().Aggregator().Clone())
		} else {
			curNode.AddAggregator(NewAggregator())
		}
	}

	for _, p := range paramValues {
		curNode.Aggregator().Add(p.Param, p.Value)
	}
}

func (t *Trie) Root() *Node {
	return t.root
}
