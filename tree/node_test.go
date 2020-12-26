package tree_test

import (
	"testing"

	"github.com/phux/vimkeypressanalyzer/tree"
	"github.com/stretchr/testify/require"
)

func TestNode_AddOrIncrementChild(t *testing.T) {
	tests := []struct {
		name         string
		existingNode *tree.Node
		expectedNode func(rootNode tree.Node) *tree.Node
		givenKeys    []string
	}{
		{
			name:         "empty after initialization",
			existingNode: tree.NewNode(""),
			expectedNode: func(rootNode tree.Node) *tree.Node { return &rootNode },
			givenKeys:    []string{},
		},
		{
			name:         "empty key",
			existingNode: tree.NewNode(""),
			expectedNode: func(rootNode tree.Node) *tree.Node { return &rootNode },
			givenKeys:    []string{""},
		},
		{
			name:         "single key",
			existingNode: tree.NewNode(""),
			expectedNode: func(rootNode tree.Node) *tree.Node {
				node := tree.NewNode("j")
				node.Count = 1
				rootNode.AddChild(node)

				return &rootNode
			},
			givenKeys: []string{"j"},
		},
		{
			name:         "3 times the same key",
			existingNode: tree.NewNode(""),
			expectedNode: func(rootNode tree.Node) *tree.Node {
				node := tree.NewNode("j")
				node.Count = 3
				rootNode.AddChild(node)

				return &rootNode
			},
			givenKeys: []string{"j", "j", "j"},
		},
		{
			name:         "alternating keys",
			existingNode: tree.NewNode(""),
			expectedNode: func(rootNode tree.Node) *tree.Node {
				nodeJ := tree.NewNode("j")
				nodeJ.Count = 2
				rootNode.AddChild(nodeJ)
				nodeI := tree.NewNode("i")
				nodeI.Count = 1
				rootNode.AddChild(nodeI)

				return &rootNode
			},
			givenKeys: []string{"j", "i", "j"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			expectedNode := tt.expectedNode(*tt.existingNode)

			for i := range tt.givenKeys {
				tt.existingNode.AddOrIncrementChild(tt.givenKeys[i])
			}

			require.Equal(t, expectedNode, tt.existingNode)
		})
	}
}
