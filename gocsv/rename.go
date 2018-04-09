package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

type RenameSubcommand struct{}

func (sub *RenameSubcommand) Name() string {
	return "rename"
}
func (sub *RenameSubcommand) Aliases() []string {
	return []string{}
}
func (sub *RenameSubcommand) Description() string {
	return "Rename the headers of a CSV."
}

func (sub *RenameSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	var columnsString, namesString string
	fs.StringVar(&columnsString, "columns", "", "Columns to rename")
	fs.StringVar(&columnsString, "c", "", "Columns to rename (shorthand)")
	fs.StringVar(&namesString, "names", "", "New names for columns")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	if columnsString == "" {
		fmt.Fprintf(os.Stderr, "Missing required argument --columns")
		os.Exit(1)
	}
	if namesString == "" {
		fmt.Fprintf(os.Stderr, "Missing required argument --names")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(columnsString)
	names := GetArrayFromCsvString(namesString)
	if len(columns) != len(names) {
		fmt.Fprintln(os.Stderr, "Length of --columns and --names argument must be the same")
		os.Exit(1)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	RenameColumns(inputCsvs[0], columns, names)
}

func RenameColumns(inputCsv AbstractInputCsv, columns, names []string) {
	writer := csv.NewWriter(os.Stdout)

	// Get the column indices to write.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	renamedHeader := make([]string, len(header))
	copy(renamedHeader, header)

	for i, column := range columns {
		index := GetColumnIndexOrPanic(header, column)
		renamedHeader[index] = names[i]
	}

	writer.Write(renamedHeader)
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
		writer.Write(row)
		writer.Flush()
	}
}
