package trie

func NewNode(key string) *Node {
	return &Node{
		key:      key,
		Children: make(map[string]*Node),
		Aggr:     nil,
	}
}

type Node struct {
	key      string
	Children map[string]*Node `json:"Children,omitempty"`
	Aggr     *Aggregator      `json:"Values,omitempty"`
}

func (n *Node) Key() string {
	return n.key
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

func (n *Node) FindNextAggregatedNode() *Node {
	nextNode := n.FirstChild()
	for !nextNode.HasAggregator() {
		nextNode = nextNode.FirstChild()
	}

	return nextNode
}
