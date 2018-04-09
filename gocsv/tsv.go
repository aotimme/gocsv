package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
)

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

func RunTsv(args []string) {
	fs := flag.NewFlagSet("tsv", flag.ExitOnError)
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
