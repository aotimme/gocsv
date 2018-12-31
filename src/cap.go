package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type CapSubcommand struct {
	namesString   string
	truncateNames bool
	defaultName   string
}

func (sub *CapSubcommand) Name() string {
	return "cap"
}
func (sub *CapSubcommand) Aliases() []string {
	return []string{}
}
func (sub *CapSubcommand) Description() string {
	return "Add a header row to a CSV."
}
func (sub *CapSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.namesString, "names", "", "Column names")
	fs.BoolVar(&sub.truncateNames, "truncate-names", false, "Truncate column names if too long")
	fs.StringVar(&sub.defaultName, "default-name", "", "Default name to use if there are more columns than column names provided")
}

func (sub *CapSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsv(inputCsvs[0])
	sub.RunCap(inputCsvs[0], outputCsv)
}

func (sub *CapSubcommand) RunCap(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	names := GetArrayFromCsvString(sub.namesString)
	Cap(inputCsv, outputCsvWriter, names, sub.truncateNames, sub.defaultName)
}

func Cap(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, names []string, truncateNames bool, defaultName string) {
	firstRow, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	numColumns := len(firstRow)
	numNames := len(names)
	if numColumns > numNames && defaultName == "" {
		fmt.Fprintf(os.Stderr, "Must specify --default-name if there are more columns than column names provided")
		os.Exit(1)
	}
	if numColumns < numNames && !truncateNames {
		fmt.Fprintf(os.Stderr, "Must specify --truncate-names if there are fewer columns than column names provided")
		os.Exit(1)
	}

	newHeader := make([]string, numColumns)
	j := 0
	for i := range firstRow {
		if i < numNames {
			newHeader[i] = names[i]
		} else {
			if j == 0 {
				newHeader[i] = defaultName
			} else {
				newHeader[i] = fmt.Sprintf("%s %d", defaultName, j)
			}
			j++
		}
	}

	outputCsvWriter.Write(newHeader)
	outputCsvWriter.Write(firstRow)

	// Write the rest of the rows.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		outputCsvWriter.Write(row)
	}
}
