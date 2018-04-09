package main

import (
	"flag"
	"fmt"
)

func ShowHeaders(inputCsv AbstractInputCsv) {
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	for i, name := range header {
		fmt.Printf("%d: %s\n", i+1, name)
	}
}

func RunHeaders(args []string) {
	fs := flag.NewFlagSet("headers", flag.ExitOnError)
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
