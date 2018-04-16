package main

import (
	"encoding/csv"
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
	numRowsRegex := regexp.MustCompile("^\\+?\\d+$")
	if !numRowsRegex.MatchString(sub.numRowsStr) {
		fmt.Fprintln(os.Stderr, "Invalid argument to -n")
		os.Exit(1)
	}

	inputCsvs, err := GetInputCsvs(args, 1)
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(sub.numRowsStr, "+") {
		numRowsStr := strings.TrimPrefix(sub.numRowsStr, "+")
		numRows, err := strconv.Atoi(numRowsStr)
		if err != nil {
			panic(err)
		}
		TailFromTop(inputCsvs[0], numRows)
	} else {
		numRows, err := strconv.Atoi(sub.numRowsStr)
		if err != nil {
			panic(err)
		}
		TailFromBottom(inputCsvs[0], numRows)
	}
}

func TailFromBottom(inputCsv AbstractInputCsv, numRows int) {
	writer := csv.NewWriter(os.Stdout)

	// Read all rows.
	rows, err := inputCsv.ReadAll()
	if err != nil {
		panic(err)
	}

	// Write header.
	writer.Write(rows[0])
	writer.Flush()

	// Write rows.
	startRow := len(rows) - numRows
	if startRow < 1 {
		startRow = 1
	}
	for i := startRow; i < len(rows); i++ {
		writer.Write(rows[i])
		writer.Flush()
	}
}

func TailFromTop(inputCsv AbstractInputCsv, numRows int) {
	writer := csv.NewWriter(os.Stdout)

	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	writer.Write(header)
	writer.Flush()

	// Write rows after first `numRows` rows.
	curRow := 0
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		curRow++
		if curRow > numRows {
			writer.Write(row)
			writer.Flush()
		}
	}
}
