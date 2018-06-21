package main

import (
	"flag"
	"fmt"
)

type HeadersSubcommand struct{}

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
}

func (sub *HeadersSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	ShowHeaders(inputCsvs[0])
}

func ShowHeaders(inputCsv AbstractInputCsv) {
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	for i, name := range header {
		fmt.Printf("%d: %s\n", i+1, name)
	}
}
