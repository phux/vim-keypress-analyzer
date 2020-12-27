package parser

import (
	"sort"

	"github.com/phux/vimkeypressanalyzer/tree"
)

type Result struct {
	KeyMap       *tree.Node
	ModeCount    *tree.Node
	Antipatterns map[string]*Antipattern
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
	total := r.TotalKeypresses([]string{})

	for i := range rootChilds {
		rootChilds[i].Share = float64(rootChilds[i].Count*100) / float64(total)
	}

	sort.Slice(rootChilds, func(i, j int) bool {
		return rootChilds[i].Count > rootChilds[j].Count
	})

	return rootChilds
}

func (r Result) SortedKeyMap(excludeModes []string) []*tree.Node {
	rootChilds := r.KeyMap.Children()
	total := r.TotalKeypresses(excludeModes)

	for i := range rootChilds {
		rootChilds[i].Share = float64(rootChilds[i].Count*100) / float64(total)
	}

	sort.Slice(rootChilds, func(i, j int) bool {
		return rootChilds[i].Count > rootChilds[j].Count
	})

	return rootChilds
}

func (r Result) TotalKeypresses(excludeModes []string) int64 {
	var total int64

NextChild:
	for _, child := range r.ModeCount.Children() {
		for _, excludeMode := range excludeModes {
			if child.Identifier == excludeMode {
				continue NextChild
			}
		}

		total += int64(child.Count)
	}

	return total
}

func (r Result) SortedAntipatterns() []*Antipattern {
	sorted := []*Antipattern{}

	for _, repetition := range r.Antipatterns {
		repetition.AverageKeypresses = float64(repetition.TotalKeypresses) / float64(repetition.Count)
		sorted = append(sorted, repetition)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	return sorted
}
