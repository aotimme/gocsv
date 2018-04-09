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

func (sub *HeadersSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}
	ShowHeaders(inputCsvs[0])
}

func ShowHeaders(inputCsv AbstractInputCsv) {
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	for i, name := range header {
		fmt.Printf("%d: %s\n", i+1, name)
	}
}
