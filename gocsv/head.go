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

func HeadFromBottom(inputCsv AbstractInputCsv, numRows int) {
	writer := csv.NewWriter(os.Stdout)

	rows, err := inputCsv.ReadAll()
	if err != nil {
		panic(err)
	}

	// Write header.
	writer.Write(rows[0])
	writer.Flush()

	// Write rows up to last `numRows` rows.
	maxRow := len(rows) - numRows
	if maxRow < 1 {
		return
	}
	for i := 1; i < maxRow; i++ {
		writer.Write(rows[i])
		writer.Flush()
	}
}

func HeadFromTop(inputCsv AbstractInputCsv, numRows int) {
	writer := csv.NewWriter(os.Stdout)

	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	writer.Write(header)
	writer.Flush()

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
				panic(err)
			}
		}
		curRow++
		writer.Write(row)
		writer.Flush()
	}
}

func RunHead(args []string) {
	fs := flag.NewFlagSet("filter", flag.ExitOnError)
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
		return
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
		HeadFromBottom(inputCsvs[0], numRows)
	} else {
		numRows, err := strconv.Atoi(numRowsStr)
		if err != nil {
			panic(err)
		}
		HeadFromTop(inputCsvs[0], numRows)
	}
}
