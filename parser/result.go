package parser

import (
	"sort"

	"github.com/phux/vimkeypressanalyzer/tree"
)

type Result struct {
	KeyMap    *tree.Node
	ModeCount *tree.Node
}

func NewResult() *Result {
	return &Result{
		KeyMap:    tree.NewNode(""),
		ModeCount: tree.NewNode(""),
	}
}

type KeyResult struct {
	Count int
}

func (kr *KeyResult) IncrCount() {
	kr.Count++
}

type SortedEntry struct {
	Name  string  `header:"name"`
	Count int     `header:"count"`
	Share float64 `header:"share (%)"`
}

func (r Result) SortedModeCount() []*tree.Node {
	rootChilds := r.ModeCount.Children()
	total := r.TotalKeypresses(true)

	for i := range rootChilds {
		rootChilds[i].Share = float64(rootChilds[i].Count*100) / float64(total)
	}

	sort.Slice(rootChilds, func(i, j int) bool {
		return rootChilds[i].Count > rootChilds[j].Count
	})

	return rootChilds
}

func (r Result) SortedKeyMap() []*tree.Node {
	rootChilds := r.KeyMap.Children()
	total := r.TotalKeypresses(false)

	for i := range rootChilds {
		rootChilds[i].Share = float64(rootChilds[i].Count*100) / float64(total)
	}

	sort.Slice(rootChilds, func(i, j int) bool {
		return rootChilds[i].Count > rootChilds[j].Count
	})

	return rootChilds
}

func (r Result) TotalKeypresses(includeInsertMode bool) int64 {
	var total int64

	for _, child := range r.ModeCount.Children() {
		if child.Identifier == InsertMode && !includeInsertMode {
			continue
		}

		total += int64(child.Count)
	}

	return total
}
