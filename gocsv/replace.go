package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Get indices to compare against.
// If no columns are specified, then check against all.
func getColumnIndicesToCompareAgainst(header, columns []string) []int {
	var columnIndices []int
	if len(columns) == 0 {
		columnIndices = make([]int, len(header))
		for i, _ := range header {
			columnIndices[i] = i
		}
	} else {
		columnIndices = make([]int, len(columns))
		for i, column := range columns {
			index := GetColumnIndexOrPanic(header, column)
			columnIndices[i] = index
		}
	}
	return columnIndices
}

func ReplaceWithFunc(reader *csv.Reader, columns []string, replaceFunc func(string) string) {
	writer := csv.NewWriter(os.Stdout)

	// Read header to get column index and write.
	header, err := reader.Read()
	if err != nil {
		panic(err)
	}

	columnIndices := getColumnIndicesToCompareAgainst(header, columns)

	writer.Write(header)
	writer.Flush()

	// Write replaced rows
	rowToWrite := make([]string, len(header))
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		copy(rowToWrite, row)
		for _, columnIndex := range columnIndices {
			rowToWrite[columnIndex] = replaceFunc(rowToWrite[columnIndex])
		}
		writer.Write(rowToWrite)
		writer.Flush()
	}
}

func RunReplace(args []string) {
	fs := flag.NewFlagSet("replace", flag.ExitOnError)
	var regex, repl, columnsString string
	var caseInsensitive bool
	fs.StringVar(&columnsString, "columns", "", "Columns to replace cells")
	fs.StringVar(&columnsString, "c", "", "Columns to replace cells (shorthand)")
	fs.StringVar(&regex, "regex", "", "Regular expression to match for replacement")
	fs.StringVar(&repl, "repl", "", "Replacement string")
	fs.BoolVar(&caseInsensitive, "case-insensitive", false, "Make regex case insensitive")
	fs.BoolVar(&caseInsensitive, "i", false, "Make regex case insensitive (shorthand)")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	// Get columns to compare against
	var columns []string
	if columnsString == "" {
		columns = make([]string, 0)
	} else {
		columns = GetArrayFromCsvString(columnsString)
	}

	// Get replace function
	var replaceFunc func(string) string
	if caseInsensitive {
		regex = "(?i)" + regex
	}
	re, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}
	replaceFunc = func(elem string) string {
		return re.ReplaceAllString(elem, repl)
	}

	// Get input CSV
	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only run replace on one table")
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

	ReplaceWithFunc(reader, columns, replaceFunc)
}
