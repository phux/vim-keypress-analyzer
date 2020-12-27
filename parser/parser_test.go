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
	type args struct {
		input io.Reader
	}

	tests := []struct {
		args              args
		expectedKeyMap    func() *tree.Node
		expectedModeCount func() *tree.Node
		name              string
		wantErr           bool
	}{
		{
			name: "empty key log",
			expectedKeyMap: func() *tree.Node {
				return tree.NewNode("")
			},
			expectedModeCount: func() *tree.Node {
				return tree.NewNode("")
			},
			args: args{
				input: strings.NewReader(""),
			},
			wantErr: false,
		},
		{
			name: "single key",
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
			args: args{
				input: strings.NewReader("j"),
			},
			wantErr: false,
		},
		{
			name: "repeating the same key",
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
			args: args{
				input: strings.NewReader("jj"),
			},
			wantErr: false,
		},
		{
			name: "escape key",
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
			args: args{
				input: strings.NewReader(string(parser.CharEsc)),
			},
			wantErr: false,
		},
		{
			name: "i<esc>",
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
			args: args{
				input: strings.NewReader("i" + string(parser.CharEsc)),
			},
			wantErr: false,
		},
		{
			name: "ji<esc>j",
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
			args: args{
				input: strings.NewReader("ji" + string(parser.CharEsc) + "j"),
			},
			wantErr: false,
		},
		{
			name: "going twice into command mode",
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
			args: args{
				input: strings.NewReader("::"),
			},
			wantErr: false,
		},
		{
			name: "going into insert mode via cc",
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
			args: args{
				input: strings.NewReader("ccc"),
			},
			wantErr: false,
		},
		{
			name: "going into insert mode via C",
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
			args: args{
				input: strings.NewReader("C" + string(parser.CharEsc) + "c"),
			},
			wantErr: false,
		},
		{
			name: "visual mode",
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
			args: args{
				input: strings.NewReader("Vj" + string(parser.CharEsc) + "vG"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := parser.NewParser()
			got, err := p.Parse(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			expectedResult := &parser.Result{
				KeyMap:    tt.expectedKeyMap(),
				ModeCount: tt.expectedModeCount(),
			}
			require.EqualValues(t, expectedResult, got)
		})
	}
}

func addChildWithCount(rootNode *tree.Node, identifier string, count int) {
	node := tree.NewNode(identifier)
	node.Count = count
	rootNode.AddChild(node)
}
