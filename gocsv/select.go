package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

func ExcludeColumns(inputCsv AbstractInputCsv, columns []string) {
	writer := csv.NewWriter(os.Stdout)

	// Get the column indices to exclude.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	columnIndicesToExclude := make(map[int]bool)
	for _, column := range columns {
		index := GetColumnIndexOrPanic(header, column)
		columnIndicesToExclude[index] = true
	}

	outrow := make([]string, len(header)-len(columnIndicesToExclude))

	// Write header
	curIdx := 0
	for index, elem := range header {
		_, exclude := columnIndicesToExclude[index]
		if !exclude {
			outrow[curIdx] = elem
			curIdx++
		}
	}

	writer.Write(outrow)
	writer.Flush()

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		curIdx = 0
		for index, elem := range row {
			_, exclude := columnIndicesToExclude[index]
			if !exclude {
				outrow[curIdx] = elem
				curIdx++
			}
		}
		writer.Write(outrow)
		writer.Flush()
	}
}

func SelectColumns(inputCsv AbstractInputCsv, columns []string) {
	writer := csv.NewWriter(os.Stdout)

	outrow := make([]string, len(columns))

	// Get the column indices to write.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	columnIndices := make([]int, len(columns))
	for i, column := range columns {
		index := GetColumnIndexOrPanic(header, column)
		columnIndices[i] = index
		outrow[i] = header[index]
	}
	writer.Write(outrow)
	writer.Flush()

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		for i, columnIndex := range columnIndices {
			outrow[i] = row[columnIndex]
		}
		writer.Write(outrow)
		writer.Flush()
	}
}

func RunSelect(args []string) {
	fs := flag.NewFlagSet("select", flag.ExitOnError)
	var columnsString string
	var exclude bool
	fs.StringVar(&columnsString, "columns", "", "Columns to select")
	fs.StringVar(&columnsString, "c", "", "Columns to select (shorthand)")
	fs.BoolVar(&exclude, "exclude", false, "Whether to exclude the specified columns")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	if columnsString == "" {
		fmt.Fprintf(os.Stderr, "Missing required argument --columns")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(columnsString)

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	if exclude {
		ExcludeColumns(inputCsvs[0], columns)
	} else {
		SelectColumns(inputCsvs[0], columns)
	}
}
