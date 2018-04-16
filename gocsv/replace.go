package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
	"regexp"
)

type ReplaceSubcommand struct{}

func (sub *ReplaceSubcommand) Name() string {
	return "replace"
}
func (sub *ReplaceSubcommand) Aliases() []string {
	return []string{}
}
func (sub *ReplaceSubcommand) Description() string {
	return "Replace values in cells by regular expression."
}

func (sub *ReplaceSubcommand) Run(args []string) {
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

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	ReplaceWithFunc(inputCsvs[0], columns, replaceFunc)
}

func ReplaceWithFunc(inputCsv AbstractInputCsv, columns []string, replaceFunc func(string) string) {
	writer := csv.NewWriter(os.Stdout)

	// Read header to get column index and write.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	writer.Write(header)
	writer.Flush()

	// Write replaced rows
	rowToWrite := make([]string, len(header))
	for {
		row, err := inputCsv.Read()
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
