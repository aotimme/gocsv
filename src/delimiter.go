package main

import (
	"flag"
	"io"
	"unicode/utf8"
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

func ChangeDelimiter(inputCsv *InputCsv, inputDelimiter, outputDelimiter string) {
	if inputDelimiter == "\\t" {
		inputCsv.SetDelimiter('\t')
	} else if len(inputDelimiter) > 0 {
		delimiterRune, _ := utf8.DecodeRuneInString(inputDelimiter)
		inputCsv.SetDelimiter(delimiterRune)
	}
	// Be lenient when reading in the file.
	inputCsv.SetFieldsPerRecord(-1)
	inputCsv.SetLazyQuotes(true)

	outputCsv := NewOutputCsvFromInputCsv(inputCsv)
	if outputDelimiter == "\\t" {
		outputCsv.SetDelimiter('\t')
	} else if len(outputDelimiter) > 0 {
		delimiterRune, _ := utf8.DecodeRuneInString(outputDelimiter)
		outputCsv.SetDelimiter(delimiterRune)
	}

	// Write all rows with tabs.
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
