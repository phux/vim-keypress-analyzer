package tree

type Node struct {
	Identifier string `header:"identifier"`
	children   []*Node
	Count      int     `header:"count"`
	Share      float64 `header:"share (%)"`
}

func NewNode(key string) *Node {
	return &Node{
		Identifier: key,
		children:   []*Node{},
	}
}

func (n *Node) AddOrIncrementChild(key string) {
	if key == "" {
		return
	}

	index := findIndexByKey(n, key)
	if index == -1 {
		index = len(n.children)
		n.AddChild(NewNode(key))
	}

	n.children[index].Count++
}

func findIndexByKey(n *Node, key string) int {
	index := -1

	for i := range n.children {
		if n.children[i].Identifier == key {
			return i
		}
	}

	return index
}

func (n *Node) AddChild(node *Node) {
	n.children = append(n.children, node)
}

func (n Node) Children() []*Node {
	return n.children
}
