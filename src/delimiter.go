package main

import (
	"flag"
	"io"
	"os"
	"unicode/utf8"

	"./csv"
)

type DelimiterSubcommand struct {
	inputDelimiter  string
	outputDelimiter string
}

func (sub *DelimiterSubcommand) Name() string {
	return "delimiter"
}
func (sub *DelimiterSubcommand) Aliases() []string {
	return []string{"delim"}
}
func (sub *DelimiterSubcommand) Description() string {
	return "Change the delimiter being used for a CSV."
}
func (sub *DelimiterSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.inputDelimiter, "input", "", "Input delimiter")
	fs.StringVar(&sub.inputDelimiter, "i", "", "Input delimiter (shorthand)")
	fs.StringVar(&sub.outputDelimiter, "output", "", "Output delimiter")
	fs.StringVar(&sub.outputDelimiter, "o", "", "Output delimiter (shorthand)")
}

func (sub *DelimiterSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	ChangeDelimiter(inputCsvs[0], sub.inputDelimiter, sub.outputDelimiter)
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
