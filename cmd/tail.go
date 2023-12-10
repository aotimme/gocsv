package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TailSubcommand struct {
	numRowsStr string
}

func (sub *TailSubcommand) Name() string {
	return "tail"
}
func (sub *TailSubcommand) Aliases() []string {
	return []string{}
}
func (sub *TailSubcommand) Description() string {
	return "Extract the last N rows from a CSV."
}
func (sub *TailSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.numRowsStr, "n", "10", "Number of rows to include")
}

func (sub *TailSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsv(inputCsvs[0])
	sub.RunTail(inputCsvs[0], outputCsv)
}

func (sub *TailSubcommand) RunTail(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	numRowsRegex := regexp.MustCompile(`^\+?\d+$`)
	if !numRowsRegex.MatchString(sub.numRowsStr) {
		fmt.Fprintln(os.Stderr, "Invalid argument to -n")
		os.Exit(1)
	}
	if strings.HasPrefix(sub.numRowsStr, "+") {
		numRowsStr := strings.TrimPrefix(sub.numRowsStr, "+")
		numRows, err := strconv.Atoi(numRowsStr)
		if err != nil {
			ExitWithError(err)
		}
		TailFromTop(inputCsv, outputCsvWriter, numRows)
	} else {
		numRows, err := strconv.Atoi(sub.numRowsStr)
		if err != nil {
			ExitWithError(err)
		}
		TailFromBottom(inputCsv, outputCsvWriter, numRows)
	}
}

func TailFromBottom(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, numRows int) {
	// Read all rows.
	rows, err := inputCsv.ReadAll()
	if err != nil {
		ExitWithError(err)
	}

	// Write header.
	outputCsvWriter.Write(rows[0])

	// Write rows.
	startRow := len(rows) - numRows
	if startRow < 1 {
		startRow = 1
	}
	for i := startRow; i < len(rows); i++ {
		outputCsvWriter.Write(rows[i])
	}
}

func TailFromTop(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, numRows int) {
	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	outputCsvWriter.Write(header)

	// Write rows after first `numRows` rows.
	curRow := 0
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		curRow++
		if curRow > numRows {
			outputCsvWriter.Write(row)
		}
	}
}
