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
	Name  string  `header:"name"`
	Count int     `header:"count"`
	Share float64 `header:"share (%)"`
}

func (r Result) SortedModeCounts() []SortedEntry {
	ss := make([]SortedEntry, 0)

	total := r.totalKeypresses(true)

	for k, v := range r.ModeCounts {
		if v < 1 {
			continue
		}

		ss = append(ss, SortedEntry{
			Name:  k,
			Count: v,
			Share: float64(v*100) / float64(total),
		})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Count > ss[j].Count
	})

	return ss
}

func (r Result) SortedKeyMap() []*tree.Node {
	rootChilds := r.KeyMap.Children()

	total := r.totalKeypresses(false)
	for i := range rootChilds {
		rootChilds[i].Share = float64(rootChilds[i].Count*100) / float64(total)
	}

	sort.Slice(rootChilds, func(i, j int) bool {
		return rootChilds[i].Count > rootChilds[j].Count
	})

	return rootChilds
}

func (r Result) totalKeypresses(includeInsertMode bool) int64 {
	var total int64
	for mode, count := range r.ModeCounts {
		if mode == InsertMode && !includeInsertMode {
			continue
		}
		total += int64(count)
	}

	return total
}
