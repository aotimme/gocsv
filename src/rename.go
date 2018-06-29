package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type RenameSubcommand struct {
	columnsString string
	namesString   string
}

func (sub *RenameSubcommand) Name() string {
	return "rename"
}
func (sub *RenameSubcommand) Aliases() []string {
	return []string{}
}
func (sub *RenameSubcommand) Description() string {
	return "Rename the headers of a CSV."
}
func (sub *RenameSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.columnsString, "columns", "", "Columns to rename")
	fs.StringVar(&sub.columnsString, "c", "", "Columns to rename (shorthand)")
	fs.StringVar(&sub.namesString, "names", "", "New names for columns")
}

func (sub *RenameSubcommand) Run(args []string) {
	if sub.columnsString == "" {
		fmt.Fprintln(os.Stderr, "Missing required argument --columns")
		os.Exit(1)
	}
	if sub.namesString == "" {
		fmt.Fprintln(os.Stderr, "Missing required argument --names")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(sub.columnsString)
	names := GetArrayFromCsvString(sub.namesString)

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	RenameColumns(inputCsvs[0], columns, names)
}

func RenameColumns(inputCsv *InputCsv, columns, names []string) {
	outputCsv := NewOutputCsvFromInputCsv(inputCsv)

	// Get the column indices to write.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	renamedHeader := make([]string, len(header))
	copy(renamedHeader, header)

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	if len(columnIndices) != len(names) {
		fmt.Fprintln(os.Stderr, "Length of --columns and --names argument must be the same")
		os.Exit(1)
	}
	for i, columnIndex := range columnIndices {
		renamedHeader[columnIndex] = names[i]
	}

	outputCsv.Write(renamedHeader)

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		outputCsv.Write(row)
	}
}
