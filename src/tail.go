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

type TailSubcommand struct{}

func (sub *TailSubcommand) Name() string {
	return "tail"
}
func (sub *TailSubcommand) Aliases() []string {
	return []string{}
}
func (sub *TailSubcommand) Description() string {
	return "Extract the last N rows from a CSV."
}

func (sub *TailSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	var numRowsStr string
	fs.StringVar(&numRowsStr, "n", "10", "Number of rows to include")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	numRowsRegex := regexp.MustCompile("^\\+?\\d+$")
	if !numRowsRegex.MatchString(numRowsStr) {
		fmt.Fprintln(os.Stderr, "Invalid argument to -n")
		os.Exit(1)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(numRowsStr, "+") {
		numRowsStr = strings.TrimPrefix(numRowsStr, "+")
		numRows, err := strconv.Atoi(numRowsStr)
		if err != nil {
			panic(err)
		}
		TailFromTop(inputCsvs[0], numRows)
	} else {
		numRows, err := strconv.Atoi(numRowsStr)
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
