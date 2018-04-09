package main

import (
	"flag"
)

func Stats(inputCsv AbstractInputCsv) {
	imc := NewInMemoryCsvFromInputCsv(inputCsv)
	imc.PrintStats()
}

func RunStats(args []string) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	Stats(inputCsvs[0])
}
