package main

import (
	"flag"
	"fmt"
	"strconv"
)

type HeadersSubcommand struct {
	asCsv bool
}

func (sub *HeadersSubcommand) Name() string {
	return "headers"
}
func (sub *HeadersSubcommand) Aliases() []string {
	return []string{}
}
func (sub *HeadersSubcommand) Description() string {
	return "View the headers from a CSV."
}
func (sub *HeadersSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&sub.asCsv, "csv", false, "Output results as CSV")
}

func (sub *HeadersSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	ShowHeaders(inputCsvs[0], sub.asCsv)
}

func ShowHeaders(inputCsv *InputCsv, asCsv bool) {
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	if asCsv {
		outputCsv := NewOutputCsvFromInputCsv(inputCsv)
		outputCsv.Write([]string{"Column", "Name"})
		for i, name := range header {
			outputCsv.Write([]string{strconv.Itoa(i + 1), name})
		}
	} else {
		for i, name := range header {
			fmt.Printf("%d: %s\n", i+1, name)
		}
	}
}
