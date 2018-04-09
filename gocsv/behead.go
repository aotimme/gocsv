package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

func Behead(inputCsv AbstractInputCsv, numHeaders int) {
	writer := csv.NewWriter(os.Stdout)

	// Get rid of the header rows.
	for i := 0; i < numHeaders; i++ {
		_, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				// If we remove _all_ the headers, then end early.
				return
			} else {
				panic(err)
			}
		}
	}

	// Write rows.
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

func RunBehead(args []string) {
	fs := flag.NewFlagSet("behead", flag.ExitOnError)
	var numHeaders int
	fs.IntVar(&numHeaders, "n", 1, "Number of headers to remove")

	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	if numHeaders < 1 {
		fmt.Fprintln(os.Stderr, "Invalid argument -n")
		os.Exit(1)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}
	Behead(inputCsvs[0], numHeaders)
}
