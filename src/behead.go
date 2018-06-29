package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type BeheadSubcommand struct {
	numHeaders int
}

func (sub *BeheadSubcommand) Name() string {
	return "behead"
}
func (sub *BeheadSubcommand) Aliases() []string {
	return []string{}
}
func (sub *BeheadSubcommand) Description() string {
	return "Remove header row(s) from a CSV."
}
func (sub *BeheadSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.IntVar(&sub.numHeaders, "n", 1, "Number of headers to remove")
}

func (sub *BeheadSubcommand) Run(args []string) {
	if sub.numHeaders < 1 {
		fmt.Fprintln(os.Stderr, "Invalid argument -n")
		os.Exit(1)
	}

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	Behead(inputCsvs[0], sub.numHeaders)
}

func Behead(inputCsv *InputCsv, numHeaders int) {
	outputCsv := NewOutputCsvFromInputCsv(inputCsv)

	// Get rid of the header rows.
	for i := 0; i < numHeaders; i++ {
		_, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				// If we remove _all_ the headers, then end early.
				return
			} else {
				ExitWithError(err)
			}
		}
	}

	// Write rows.
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
