package main

import (
	"flag"
	"io"
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

func Tsv(inputCsv *InputCsv) {
	outputCsv := NewOutputCsvFromInputCsv(inputCsv)
	outputCsv.SetDelimiter('\t')

	// Write all rows with tabs.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		outputCsv.Write(row)
	}
}
