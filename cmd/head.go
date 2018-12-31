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

type HeadSubcommand struct {
	numRowsStr string
}

func (sub *HeadSubcommand) Name() string {
	return "head"
}
func (sub *HeadSubcommand) Aliases() []string {
	return []string{}
}
func (sub *HeadSubcommand) Description() string {
	return "Extract the first N rows from a CSV."
}
func (sub *HeadSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.numRowsStr, "n", "10", "Number of rows to include")
}

func (sub *HeadSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsv(inputCsvs[0])
	sub.RunHead(inputCsvs[0], outputCsv)
}

func (sub *HeadSubcommand) RunHead(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	numRowsRegex := regexp.MustCompile("^\\+?\\d+$")
	if !numRowsRegex.MatchString(sub.numRowsStr) {
		fmt.Fprintln(os.Stderr, "Invalid argument to -n")
		os.Exit(1)
		return
	}

	if strings.HasPrefix(sub.numRowsStr, "+") {
		sub.numRowsStr = strings.TrimPrefix(sub.numRowsStr, "+")
		numRows, err := strconv.Atoi(sub.numRowsStr)
		if err != nil {
			ExitWithError(err)
		}
		HeadFromBottom(inputCsv, outputCsvWriter, numRows)
	} else {
		numRows, err := strconv.Atoi(sub.numRowsStr)
		if err != nil {
			ExitWithError(err)
		}
		HeadFromTop(inputCsv, outputCsvWriter, numRows)
	}
}

func HeadFromBottom(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, numRows int) {
	rows, err := inputCsv.ReadAll()
	if err != nil {
		ExitWithError(err)
	}

	// Write header.
	outputCsvWriter.Write(rows[0])

	// Write rows up to last `numRows` rows.
	maxRow := len(rows) - numRows
	if maxRow < 1 {
		return
	}
	for i := 1; i < maxRow; i++ {
		outputCsvWriter.Write(rows[i])
	}
}

func HeadFromTop(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, numRows int) {
	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	outputCsvWriter.Write(header)

	// Write first `numRows` rows.
	curRow := 0
	for {
		if curRow == numRows {
			break
		}
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		curRow++
		outputCsvWriter.Write(row)
	}
}
