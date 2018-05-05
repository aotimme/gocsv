package main

import (
	"encoding/csv"
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
	names := GetArrayFromCsvString(sub.namesString)
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	Cap(inputCsvs[0], names, sub.truncateNames, sub.defaultName)
}

func Cap(inputCsv AbstractInputCsv, names []string, truncateNames bool, defaultName string) {
	writer := csv.NewWriter(os.Stdout)

	firstRow, err := inputCsv.Read()
	if err != nil {
		panic(err)
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

	writer.Write(newHeader)
	writer.Flush()

	writer.Write(firstRow)
	writer.Flush()

	// Write the rest of the rows.
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
