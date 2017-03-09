package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

func ExcludeColumns(reader *csv.Reader, columns []string) {
	writer := csv.NewWriter(os.Stdout)

	// Get the column indices to exclude.
	header, err := reader.Read()
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
		row, err := reader.Read()
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

func SelectColumns(reader *csv.Reader, columns []string) {
	writer := csv.NewWriter(os.Stdout)

	outrow := make([]string, len(columns))

	// Get the column indices to write.
	header, err := reader.Read()
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
		row, err := reader.Read()
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
	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only select one table")
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
	if exclude {
		ExcludeColumns(reader, columns)
	} else {
		SelectColumns(reader, columns)
	}
}
