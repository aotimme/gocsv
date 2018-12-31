package main

import (
	"flag"
	"fmt"
	"os"
)

type SortSubcommand struct {
	columnsString string
	reverse       bool
	noInference   bool
}

func (sub *SortSubcommand) Name() string {
	return "sort"
}
func (sub *SortSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SortSubcommand) Description() string {
	return "Sort a CSV based on one or more columns."
}
func (sub *SortSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.columnsString, "columns", "", "Columns to select")
	fs.StringVar(&sub.columnsString, "c", "", "Columns to select (shorthand)")
	fs.BoolVar(&sub.reverse, "reverse", false, "Sort in reverse")
	fs.BoolVar(&sub.noInference, "no-inference", false, "Skip inference of input")
}

func (sub *SortSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsvs(inputCsvs)
	sub.SortCsv(inputCsvs[0], outputCsv)
}

func (sub *SortSubcommand) SortCsv(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	if sub.columnsString == "" {
		fmt.Fprintln(os.Stderr, "Missing required argument --columns")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(sub.columnsString)

	imc := NewInMemoryCsvFromInputCsv(inputCsv)
	columnIndices := GetIndicesForColumnsOrPanic(imc.header, columns)
	columnTypes := make([]ColumnType, len(columnIndices))
	for i, columnIndex := range columnIndices {
		if sub.noInference {
			columnTypes[i] = STRING_TYPE
		} else {
			columnTypes[i] = imc.InferType(columnIndex)
		}
	}
	imc.SortRows(columnIndices, columnTypes, sub.reverse)

	// Write header.
	outputCsvWriter.Write(imc.header)

	// Write sorted rows.
	for _, row := range imc.rows {
		outputCsvWriter.Write(row)
	}
}
