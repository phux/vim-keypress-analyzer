package parser

import (
	"sort"

	"github.com/phux/vimkeypressanalyzer/tree"
)

type Result struct {
	KeyMap     *tree.Node
	ModeCounts map[string]int
}

func NewResult() *Result {
	return &Result{
		KeyMap:     tree.NewNode(""),
		ModeCounts: map[string]int{},
	}
}

type KeyResult struct {
	Count int
}

func (kr *KeyResult) IncrCount() {
	kr.Count++
}

func (r *Result) IncrModeCount(mode string) {
	if _, ok := r.ModeCounts[mode]; !ok {
		r.ModeCounts[mode] = 0
	}
	r.ModeCounts[mode]++
}

type SortedEntry struct {
	Name  string `header:"name"`
	Count int    `header:"count"`
}

func (r Result) SortedModeCounts() []SortedEntry {
	ss := make([]SortedEntry, 0)

	for k, v := range r.ModeCounts {
		if v < 1 {
			continue
		}

		ss = append(ss, SortedEntry{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Count > ss[j].Count
	})

	return ss
}

func (r Result) SortedKeyMap() []*tree.Node {
	rootChilds := r.KeyMap.Children()

	sort.Slice(rootChilds, func(i, j int) bool {
		return rootChilds[i].Count > rootChilds[j].Count
	})

	return rootChilds
}
