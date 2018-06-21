package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"./csv"
)

type SelectSubcommand struct {
	columnsString string
	exclude       bool
}

func (sub *SelectSubcommand) Name() string {
	return "select"
}
func (sub *SelectSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SelectSubcommand) Description() string {
	return "Extract specified columns."
}
func (sub *SelectSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.columnsString, "columns", "", "Columns to select")
	fs.StringVar(&sub.columnsString, "c", "", "Columns to select (shorthand)")
	fs.BoolVar(&sub.exclude, "exclude", false, "Whether to exclude the specified columns")
}

func (sub *SelectSubcommand) Run(args []string) {
	if sub.columnsString == "" {
		fmt.Fprintln(os.Stderr, "Missing required argument --columns")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(sub.columnsString)

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	if sub.exclude {
		ExcludeColumns(inputCsvs[0], columns)
	} else {
		SelectColumns(inputCsvs[0], columns)
	}
}

func ExcludeColumns(inputCsv AbstractInputCsv, columns []string) {
	writer := csv.NewWriter(os.Stdout)

	// Get the column indices to exclude.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	columnIndices := GetIndicesForColumnsOrPanic(header, columns)
	columnIndicesToExclude := make(map[int]bool)
	for _, columnIndex := range columnIndices {
		columnIndicesToExclude[columnIndex] = true
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
				ExitWithError(err)
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

	// Get the column indices to write.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	columnIndices := GetIndicesForColumnsOrPanic(header, columns)
	outrow := make([]string, len(columnIndices))
	for i, columnIndex := range columnIndices {
		outrow[i] = header[columnIndex]
	}
	writer.Write(outrow)
	writer.Flush()

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		for i, columnIndex := range columnIndices {
			outrow[i] = row[columnIndex]
		}
		writer.Write(outrow)
		writer.Flush()
	}
}
