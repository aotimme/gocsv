package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func Sample(reader *csv.Reader, numRows int, replace bool, seed int) {

	imc := NewInMemoryCsv(reader)

	if numRows > imc.NumRows() && !replace {
		fmt.Fprintln(os.Stderr, "Cannot sample more rows than exist")
		os.Exit(1)
	}

	rowIndices := imc.SampleRowIndices(numRows, replace, seed)

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	writer.Write(imc.header)
	writer.Flush()

	for _, rowIndex := range rowIndices {
		writer.Write(imc.rows[rowIndex])
		writer.Flush()
	}

}

func RunSample(args []string) {
	fs := flag.NewFlagSet("sample", flag.ExitOnError)
	var replace bool
	var numRows, seed int
	fs.BoolVar(&replace, "replace", false, "Sample with replacement")
	fs.IntVar(&numRows, "n", 0, "Number of rows to sample")
	fs.IntVar(&seed, "seed", 0, "Seed for random number generation")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	if numRows < 1 {
		fmt.Fprintln(os.Stderr, "Invalid required argument -n")
		os.Exit(1)
	}

	// Get input CSV
	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only sample one table")
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

	Sample(reader, numRows, replace, seed)
}
