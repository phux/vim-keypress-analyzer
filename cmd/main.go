package main

import (
	"log"
	"os"

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
	}
	cliApp := &cli.App{
		Flags: flags,
		Name:  "vimkeypressanalyzer",
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
			parser := parser.NewParser()
			a := app.NewApp(parser)

			return a.Analyze(logContents, c.Int64("limit"))
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
