package tree

type Node struct {
	Key      string `header:"key"`
	children []*Node
	Count    int `header:"count"`
}

func NewNode(key string) *Node {
	return &Node{
		Key:      key,
		children: []*Node{},
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
		if n.children[i].Key == key {
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
