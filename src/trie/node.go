package trie

func NewNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
		Aggr:     nil,
	}
}

type Node struct {
	Children map[string]*Node `json:"Children,omitempty"`
	Aggr     *Aggregator      `json:"Values,omitempty"`
}

func (n *Node) HasAggregator() bool {
	return n.Aggr != nil
}

func (n *Node) HasChild(key string) bool {
	return n.Children[key] != nil
}

func (n *Node) GetChild(key string) *Node {
	return n.Children[key]
}

func (n *Node) Aggregator() *Aggregator {
	return n.Aggr
}

func (n *Node) IsFork() bool {
	return len(n.Children) > 1
}

func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}

func (n *Node) AddAggregator(aggr *Aggregator) {
	n.Aggr = aggr
}

func (n *Node) AddChild(key string, node *Node) {
	n.Children[key] = node
}

func (n *Node) ChildCount() int {
	return len(n.Children)
}

func (n *Node) FirstChild() *Node {
	for _, child := range n.Children {
		return child
	}

	panic("No Children when called FirstChild")
}

func (n *Node) FirstChildWithKey() (string, *Node) {
	for key, child := range n.Children {
		return key, child
	}

	panic("No Children when called FirstChildWithKey")
}

func (n *Node) FindNextAggregatedNode() *Node {
	nextNode := n.FirstChild()
	for !nextNode.HasAggregator() {
		nextNode = nextNode.FirstChild()
	}

	return nextNode
}

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

//func (n *Node) FastForward() ([]string, *Node) {
//	if false == n.HasChildren() {
//		return []string{}, nil
//	}
//
//	key, nextNode := n.FirstChild()
//	keys := []string{key}
//	for !nextNode.HasAggregator() && n.HasChildren() {
//		var k string
//		k, nextNode = nextNode.FirstChild()
//		key = key + k
//	}
//
//	return key, nextNode
//}
