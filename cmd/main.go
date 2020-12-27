package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lensesio/tableprinter"
	"github.com/phux/vimkeypressanalyzer/app"
	"github.com/phux/vimkeypressanalyzer/parser"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Value:    "",
			Aliases:  []string{"f"},
			Usage:    "path to logfile containing the keystrokes",
			Required: true,
		},
		&cli.Int64Flag{
			Name:    "limit",
			Aliases: []string{"l"},
			Usage:   "number of most frequent keys to show",
			Value:   0,
		},
		&cli.BoolFlag{
			Name:    "enable-antipatterns",
			Aliases: []string{"a"},
			Usage:   "enable naive antipattern analysis",
			Value:   false,
		},
	}
	cliApp := &cli.App{
		Flags: flags,
		Name:  "Vim Keypress Analyzer",
		Usage: "parse the pressed keys in vim and give a helpful analysis",
		Action: func(c *cli.Context) error {
			logfile := c.String("file")
			if logfile == "" {
				return errors.New("no logfile given")
			}
			logContents, err := os.Open(logfile)
			if err != nil {
				return errors.Wrapf(err, "could not open logfile '%s'", logfile)
			}

			enableAntipatterns := c.Bool("enable-antipatterns")
			parser := parser.NewParser(enableAntipatterns)
			a := app.NewApp(parser)

			result, err := a.Analyze(logContents, c.Int64("limit"))
			if err != nil {
				return errors.Wrapf(err, "cmd: failed to analyze %s", logfile)
			}

			fmt.Printf("\nVim Keypress Analyzer\n\n")

			printer := tableprinter.New(os.Stdout)
			printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
			printer.CenterSeparator = "│"
			printer.ColumnSeparator = "│"
			printer.RowSeparator = "─"

			fmt.Printf("Key presses per mode (total: %d)\n", result.TotalKeypresses)

			printer.Print(result.SortedModeCounts)

			fmt.Printf(
				"\nKey presses in non-INSERT modes (total: %d)\n",
				result.TotalKeypressesWithoutInsertMode,
			)

			printer.Print(result.SortedKeyMap)

			if enableAntipatterns {
				fmt.Println("\nAntipatterns (naive approach)")

				if len(result.SortedAntipatterns) == 0 {
					fmt.Println("no antipatterns found, good job :)")
				}
				printer.Print(result.SortedAntipatterns)
			}

			return nil
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
