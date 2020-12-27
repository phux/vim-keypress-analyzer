package app

import (
	"io"

	"github.com/phux/vimkeypressanalyzer/parser"
	"github.com/phux/vimkeypressanalyzer/tree"
	"github.com/pkg/errors"
)

type Parser interface {
	Parse(log io.Reader) (*parser.Result, error)
}

type App struct {
	parser Parser
}

type AnalyzerResult struct {
	SortedModeCounts                 []*tree.Node
	SortedKeyMap                     []*tree.Node
	TotalKeypresses                  int64
	TotalKeypressesWithoutInsertMode int64
}

func NewApp(p Parser) App {
	return App{
		parser: p,
	}
}

func (a App) Analyze(log io.Reader, limit int64) (AnalyzerResult, error) {
	analyzerResult := AnalyzerResult{}

	parserResult, err := a.parser.Parse(log)
	if err != nil {
		return analyzerResult, errors.Wrap(err, "could not analyze vim log")
	}

	analyzerResult.SortedModeCounts = parserResult.SortedModeCount()
	analyzerResult.TotalKeypresses = parserResult.TotalKeypresses(true)

	analyzerResult.TotalKeypressesWithoutInsertMode = parserResult.TotalKeypresses(false)

	sortedKeyMap := parserResult.SortedKeyMap()
	if limit > 0 && len(sortedKeyMap) > int(limit) {
		sortedKeyMap = sortedKeyMap[0:limit]
	}

	analyzerResult.SortedKeyMap = sortedKeyMap

	return analyzerResult, nil
}
