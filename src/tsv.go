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
func (sub *TsvSubcommand) SetFlags(fs *flag.FlagSet) {
}

func (sub *TsvSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
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
