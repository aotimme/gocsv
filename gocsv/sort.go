package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

type SortSubcommand struct{}

func (sub *SortSubcommand) Name() string {
	return "sort"
}
func (sub *SortSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SortSubcommand) Description() string {
	return "Sort a CSV based on one or more columns."
}

func (sub *SortSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	var columnsString string
	var reverse, noInference bool
	fs.StringVar(&columnsString, "columns", "", "Columns to sort on")
	fs.StringVar(&columnsString, "c", "", "Columns to sort on (shorthand)")
	fs.BoolVar(&reverse, "reverse", false, "Sort in reverse")
	fs.BoolVar(&noInference, "no-inference", false, "Skip inference of input")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	if columnsString == "" {
		fmt.Fprintln(os.Stderr, "Missing argument --columns")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(columnsString)

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	SortCsv(inputCsvs[0], columns, reverse, noInference)
}

func SortCsv(inputCsv AbstractInputCsv, columns []string, reverse, noInference bool) {
	imc := NewInMemoryCsvFromInputCsv(inputCsv)
	columnIndices := GetIndicesForColumnsOrPanic(imc.header, columns)
	columnTypes := make([]ColumnType, len(columnIndices))
	for i, columnIndex := range columnIndices {
		if noInference {
			columnTypes[i] = STRING_TYPE
		} else {
			columnTypes[i] = imc.InferType(columnIndex)
		}
	}
	imc.SortRows(columnIndices, columnTypes, reverse)

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	writer.Write(imc.header)
	writer.Flush()

	// Write sorted rows.
	for _, row := range imc.rows {
		writer.Write(row)
		writer.Flush()
	}
}
