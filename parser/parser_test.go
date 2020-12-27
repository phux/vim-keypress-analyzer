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
		want              *parser.Result
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
				node := tree.NewNode("j")
				node.Count = 1
				rootNode.AddChild(node)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				node := tree.NewNode(parser.NormalMode)
				node.Count = 1
				rootNode.AddChild(node)

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
				node := tree.NewNode("j")
				node.Count = 2
				rootNode.AddChild(node)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				nodeNormal := tree.NewNode(parser.NormalMode)
				nodeNormal.Count = 2
				rootNode.AddChild(nodeNormal)

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
				node := tree.NewNode(parser.CharReadableEsc)
				node.Count = 1
				rootNode.AddChild(node)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				nodeNormal := tree.NewNode(parser.NormalMode)
				nodeNormal.Count = 1
				rootNode.AddChild(nodeNormal)

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
				nodeI := tree.NewNode("i")
				nodeI.Count = 1
				rootNode.AddChild(nodeI)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				nodeNormal := tree.NewNode(parser.NormalMode)
				nodeNormal.Count = 1
				rootNode.AddChild(nodeNormal)
				nodeInsert := tree.NewNode(parser.InsertMode)
				nodeInsert.Count = 1
				rootNode.AddChild(nodeInsert)

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
				nodeJ := tree.NewNode("j")
				nodeJ.Count = 2
				rootNode.AddChild(nodeJ)
				nodeI := tree.NewNode("i")
				nodeI.Count = 1
				rootNode.AddChild(nodeI)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				nodeNormal := tree.NewNode(parser.NormalMode)
				nodeNormal.Count = 3
				rootNode.AddChild(nodeNormal)
				nodeInsert := tree.NewNode(parser.InsertMode)
				nodeInsert.Count = 1
				rootNode.AddChild(nodeInsert)

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
				node := tree.NewNode(":")
				node.Count = 2
				rootNode.AddChild(node)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				nodeNormal := tree.NewNode(parser.NormalMode)
				nodeNormal.Count = 1
				rootNode.AddChild(nodeNormal)
				nodeCommand := tree.NewNode(parser.CommandMode)
				nodeCommand.Count = 1
				rootNode.AddChild(nodeCommand)

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
				node := tree.NewNode("c")
				node.Count = 2
				rootNode.AddChild(node)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				nodeNormal := tree.NewNode(parser.NormalMode)
				nodeNormal.Count = 2
				rootNode.AddChild(nodeNormal)
				nodeInsert := tree.NewNode(parser.InsertMode)
				nodeInsert.Count = 1
				rootNode.AddChild(nodeInsert)

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
				node := tree.NewNode("C")
				node.Count = 1
				rootNode.AddChild(node)
				nodeC := tree.NewNode("c")
				nodeC.Count = 1
				rootNode.AddChild(nodeC)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				nodeNormal := tree.NewNode(parser.NormalMode)
				nodeNormal.Count = 2
				rootNode.AddChild(nodeNormal)
				nodeInsert := tree.NewNode(parser.InsertMode)
				nodeInsert.Count = 1
				rootNode.AddChild(nodeInsert)

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
				nodeV := tree.NewNode("V")
				nodeV.Count = 1
				rootNode.AddChild(nodeV)
				nodej := tree.NewNode("j")
				nodej.Count = 1
				rootNode.AddChild(nodej)
				nodeEsc := tree.NewNode(parser.CharReadableEsc)
				nodeEsc.Count = 1
				rootNode.AddChild(nodeEsc)

				nodev := tree.NewNode("v")
				nodev.Count = 1
				rootNode.AddChild(nodev)
				nodeG := tree.NewNode("G")
				nodeG.Count = 1
				rootNode.AddChild(nodeG)

				return rootNode
			},
			expectedModeCount: func() *tree.Node {
				rootNode := tree.NewNode("")
				node := tree.NewNode(parser.NormalMode)
				node.Count = 2
				rootNode.AddChild(node)
				nodeVisual := tree.NewNode(parser.VisualMode)
				nodeVisual.Count = 3
				rootNode.AddChild(nodeVisual)

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
