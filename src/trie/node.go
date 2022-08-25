package trie

import "sync"

// NewNode creates Node
func NewNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
		Aggr:     nil,
	}
}

// Node node of the trie
type Node struct {
	childrenLock          sync.RWMutex
	Children              map[string]*Node `json:"Children,omitempty"`
	Aggr                  *Aggregator      `json:"Values,omitempty"`
	OverflowChildrenCount uint64           `json:"Overflow,omitempty"`
}

// HasAggregator returns if the node has aggregator attached
func (n *Node) HasAggregator() bool {
	return n.Aggr != nil
}

// GetChild returns a child node by provided key, if key doesn't exist returns nil
func (n *Node) GetChild(key string) *Node {
	n.childrenLock.RLock()
	defer n.childrenLock.RUnlock()

	return n.Children[key]
}

// Aggregator Returns aggregator attached to the node
func (n *Node) Aggregator() *Aggregator {
	return n.Aggr
}

// HasChildren returns if the node has at least one child node
func (n *Node) HasChildren() bool {
	n.childrenLock.RLock()
	defer n.childrenLock.RUnlock()

	return len(n.Children) > 0
}

// AddAggregator adds aggregator
func (n *Node) AddAggregator(aggr *Aggregator) {
	n.Aggr = aggr
}

// AddChild add child nodes on provided key
func (n *Node) AddChild(key string, node *Node) {
	n.childrenLock.Lock()
	defer n.childrenLock.Unlock()

	n.Children[key] = node
}

// ChildCount return number of child nodes
func (n *Node) ChildCount() int {
	n.childrenLock.RLock()
	defer n.childrenLock.RUnlock()

	return len(n.Children)
}

// FirstChild returns child node, panics if there are no child nodes
func (n *Node) FirstChild() *Node {
	n.childrenLock.RLock()
	defer n.childrenLock.RUnlock()

	for _, child := range n.Children {
		return child
	}

	panic("No Children when called FirstChild")
}

// FirstChildWithKey return child node and its key, panics if there are no child nodes
func (n *Node) FirstChildWithKey() (string, *Node) {
	n.childrenLock.RLock()
	defer n.childrenLock.RUnlock()

	for key, child := range n.Children {
		return key, child
	}

	panic("No Children when called FirstChildWithKey")
}

// FindNextAggregatedNode descends down trie branch till the next node with aggregator
func (n *Node) FindNextAggregatedNode() *Node {
	nextNode := n.FirstChild()
	for !nextNode.HasAggregator() {
		nextNode = nextNode.FirstChild()
	}

	return nextNode
}

// FindNextAggregatedNodeWithKey descends down trie branch till the next node with aggregator and returns key path to it
func (n *Node) FindNextAggregatedNodeWithKey() ([]string, *Node) {
	firstKey, nextNode := n.FirstChildWithKey()
	keys := []string{firstKey}
	for !nextNode.HasAggregator() {
		var key string
		key, nextNode = nextNode.FirstChildWithKey()
		keys = append(keys, key)
	}

	return keys, nextNode
}
