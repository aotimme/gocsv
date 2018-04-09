package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
)

type TsvSubcommand struct{}

func (sub *TsvSubcommand) Name() string {
	return "tsv"
}
func (sub *TsvSubcommand) Aliases() []string {
	return []string{}
}
func (sub *TsvSubcommand) Description() string {
	return "Transform a CSV into a TSV."
}

func (sub *TsvSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}
	Tsv(inputCsvs[0])
}

func Tsv(inputCsv AbstractInputCsv) {
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = '\t'

	// Write all rows with tabs.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		writer.Write(row)
		writer.Flush()
	}
}
