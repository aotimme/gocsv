package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func Stats(reader *csv.Reader) {
	imc := NewInMemoryCsv(reader)
	imc.PrintStats()
}

func RunStats(args []string) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	// Get input CSV
	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only get stats on one table")
		os.Exit(1)
	}
	var reader *csv.Reader
	if len(moreArgs) == 1 {
		file, err := os.Open(moreArgs[0])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader = csv.NewReader(file)
	} else {
		reader = csv.NewReader(os.Stdin)
	}

	Stats(reader)
}
