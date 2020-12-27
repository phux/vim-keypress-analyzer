package parser_test

import (
	"io"
	"strings"
	"testing"

	"github.com/phux/vimkeypressanalyzer/parser"
	"github.com/phux/vimkeypressanalyzer/tree"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input             io.Reader
		expectedKeyMap    func() *tree.Node
		expectedModeCount func() *tree.Node
		name              string
	}{
		{
			name:  "empty key log",
			input: strings.NewReader(""),
			expectedKeyMap: func() *tree.Node {
				return tree.NewNode("")
			},
			expectedModeCount: func() *tree.Node {
				return tree.NewNode("")
			},
		},
		{
			name:  "single key",
			input: strings.NewReader("j"),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, "j", 1)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 1)

				return rootNode
			},
		},
		{
			name:  "repeating the same key",
			input: strings.NewReader("jj"),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, "j", 2)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 2)

				return rootNode
			},
		},
		{
			name:  "escape key",
			input: strings.NewReader(string(parser.CharEsc)),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.CharReadableEsc, 1)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 1)

				return rootNode
			},
		},
		{
			name:  "i<esc>",
			input: strings.NewReader("i" + string(parser.CharEsc)),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, "i", 1)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 1)
				addChildWithCount(rootNode, parser.InsertMode, 1)

				return rootNode
			},
		},
		{
			name:  "ji<esc>j",
			input: strings.NewReader("ji" + string(parser.CharEsc) + "j"),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, "j", 2)
				addChildWithCount(rootNode, "i", 1)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 3)
				addChildWithCount(rootNode, parser.InsertMode, 1)

				return rootNode
			},
		},
		{
			name:  "going twice into command mode",
			input: strings.NewReader("::"),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, ":", 2)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 1)
				addChildWithCount(rootNode, parser.CommandMode, 1)

				return rootNode
			},
		},
		{
			name:  "going into insert mode via cc",
			input: strings.NewReader("ccc"),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, "c", 2)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 2)
				addChildWithCount(rootNode, parser.InsertMode, 1)

				return rootNode
			},
		},
		{
			name:  "going into insert mode via C",
			input: strings.NewReader("C" + string(parser.CharEsc) + "c"),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, "C", 1)
				addChildWithCount(rootNode, "c", 1)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 2)
				addChildWithCount(rootNode, parser.InsertMode, 1)

				return rootNode
			},
		},
		{
			name:  "visual mode",
			input: strings.NewReader("Vj" + string(parser.CharEsc) + "vG"),
			expectedKeyMap: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, "V", 1)
				addChildWithCount(rootNode, "j", 1)
				addChildWithCount(rootNode, parser.CharReadableEsc, 1)
				addChildWithCount(rootNode, "v", 1)
				addChildWithCount(rootNode, "G", 1)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				addChildWithCount(rootNode, parser.NormalMode, 2)
				addChildWithCount(rootNode, parser.VisualMode, 3)

				return rootNode
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := parser.NewParser()

			got, err := p.Parse(tt.input)
			require.NoError(t, err)

			require.EqualValues(t, tt.expectedKeyMap(), got.KeyMap)
			require.EqualValues(t, tt.expectedModeCount(), got.ModeCount)
		})
	}
}

func addChildWithCount(rootNode *tree.Node, identifier string, count int) {
	node := tree.NewNode(identifier)
	node.Count = count
	rootNode.AddChild(node)
}
