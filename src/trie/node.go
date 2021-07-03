package trie

func NewNode(key string) *Node {
	return &Node{
		key:        key,
		children:   make(map[string]*Node),
		aggregator: nil,
	}
}

type Node struct {
	key        string
	children   map[string]*Node
	aggregator *Aggregator
}

func (n *Node) Key() string {
	return n.key
}

func (n *Node) HasAggregator() bool {
	return n.aggregator != nil
}

func (n *Node) HasChild(key string) bool {
	return n.children[key] != nil
}

func (n *Node) GetChild(key string) *Node {
	return n.children[key]
}

func (n *Node) Aggregator() *Aggregator {
	return n.aggregator
}

func (n *Node) IsFork() bool {
	return len(n.children) > 1
}

func (n *Node) AddAggregator(aggregator *Aggregator) {
	n.aggregator = aggregator
}

func (n *Node) AddChild(key string, node *Node) {
	n.children[key] = node
}
