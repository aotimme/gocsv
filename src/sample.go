package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

type SampleSubcommand struct{}

func (sub *SampleSubcommand) Name() string {
	return "sample"
}
func (sub *SampleSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SampleSubcommand) Description() string {
	return "Sample rows."
}

func (sub *SampleSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
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

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	Sample(inputCsvs[0], numRows, replace, seed)
}

func Sample(inputCsv AbstractInputCsv, numRows int, replace bool, seed int) {

	imc := NewInMemoryCsvFromInputCsv(inputCsv)

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
