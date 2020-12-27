package app

import (
	"io"

	"github.com/phux/vimkeypressanalyzer/parser"
	"github.com/phux/vimkeypressanalyzer/tree"
	"github.com/pkg/errors"
)

type Parser interface {
	Parse(log io.Reader, excludeModes []string) (*parser.Result, error)
}

type App struct {
	parser Parser
}

type AnalyzerResult struct {
	SortedModeCounts                    []*tree.Node
	SortedKeyMap                        []*tree.Node
	SortedAntipatterns                  []*parser.Antipattern
	TotalKeypresses                     int64
	TotalKeypressesWithoutExcludedModes int64
}

func NewApp(p Parser) App {
	return App{
		parser: p,
	}
}

func (a App) Analyze(log io.Reader, limit int64, excludeModes []string) (AnalyzerResult, error) {
	analyzerResult := AnalyzerResult{}

	for _, mode := range excludeModes {
		switch mode {
		case parser.NormalMode, parser.InsertMode, parser.CommandMode, parser.VisualMode:
		default:
			return analyzerResult, errors.Errorf("invalid exclude-mode given: %s", mode)
		}
	}

	parserResult, err := a.parser.Parse(log, excludeModes)
	if err != nil {
		return analyzerResult, errors.Wrap(err, "could not analyze vim log")
	}

	analyzerResult.SortedModeCounts = parserResult.SortedModeCount()
	analyzerResult.TotalKeypresses = parserResult.TotalKeypresses(nil)

	analyzerResult.TotalKeypressesWithoutExcludedModes = parserResult.TotalKeypresses(excludeModes)

	sortedKeyMap := parserResult.SortedKeyMap(excludeModes)
	if limit > 0 && len(sortedKeyMap) > int(limit) {
		sortedKeyMap = sortedKeyMap[0:limit]
	}

	analyzerResult.SortedKeyMap = sortedKeyMap

	analyzerResult.SortedAntipatterns = parserResult.SortedAntipatterns()

	return analyzerResult, nil
}
