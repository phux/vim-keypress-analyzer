package app

import (
	"fmt"
	"io"
	"os"

	"github.com/lensesio/tableprinter"
	"github.com/phux/vimkeypressanalyzer/parser"
	"github.com/pkg/errors"
)

type Parser interface {
	Parse(log io.Reader) (*parser.Result, error)
}

type App struct {
	parser Parser
}

func NewApp(p Parser) App {
	return App{
		parser: p,
	}
}

func (a App) Analyze(log io.Reader, limit int64) error {
	result, err := a.parser.Parse(log)
	if err != nil {
		return errors.Wrap(err, "could not analyze vim log")
	}

	fmt.Printf("\nVim Keypress Analyzer\n\n")

	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"

	sortedModeCounts := result.SortedModeCounts()

	totalSum := 0
	for _, mode := range sortedModeCounts {
		totalSum += mode.Count
	}

	fmt.Printf("Key presses per mode (total: %d)\n", totalSum)

	printer.Print(sortedModeCounts)

	fmt.Println("\nKey presses (in non-INSERT mode)")

	sortedKeys := result.SortedKeyMap()
	if limit > 0 {
		sortedKeys = sortedKeys[0:limit]
	}

	printer.Print(sortedKeys)

	return nil
}
