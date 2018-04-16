package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
	"unicode/utf8"
)

type DelimiterSubcommand struct{}

func (sub *DelimiterSubcommand) Name() string {
	return "delimiter"
}
func (sub *DelimiterSubcommand) Aliases() []string {
	return []string{"delim"}
}
func (sub *DelimiterSubcommand) Description() string {
	return "Change the delimiter being used for a CSV."
}

func (sub *DelimiterSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	var inputDelimiter, outputDelimiter string
	fs.StringVar(&inputDelimiter, "input", "", "Input delimiter")
	fs.StringVar(&inputDelimiter, "i", "", "Input delimiter (shorthand)")
	fs.StringVar(&outputDelimiter, "output", "", "Output delimiter")
	fs.StringVar(&outputDelimiter, "o", "", "Output delimiter (shorthand)")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}
	ChangeDelimiter(inputCsvs[0], inputDelimiter, outputDelimiter)
}

func ChangeDelimiter(inputCsv AbstractInputCsv, inputDelimiter, outputDelimiter string) {
	reader := inputCsv.Reader()
	if inputDelimiter == "\\t" {
		reader.Comma = '\t'
	} else if len(inputDelimiter) > 0 {
		reader.Comma, _ = utf8.DecodeRuneInString(inputDelimiter)
	}
	// Be lenient when reading in the file.
	reader.FieldsPerRecord = -1

	writer := csv.NewWriter(os.Stdout)
	if outputDelimiter == "\\t" {
		writer.Comma = '\t'
	} else if len(outputDelimiter) > 0 {
		writer.Comma, _ = utf8.DecodeRuneInString(outputDelimiter)
	}

	// Write all rows with tabs.
	for {
		row, err := reader.Read()
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
