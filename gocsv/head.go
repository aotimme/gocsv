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

func HeadFromBottom(reader *csv.Reader, numRows int) {
	writer := csv.NewWriter(os.Stdout)

	rows, err := reader.ReadAll()
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

func HeadFromTop(reader *csv.Reader, numRows int) {
	writer := csv.NewWriter(os.Stdout)

	// Read and write header.
	header, err := reader.Read()
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
		row, err := reader.Read()
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
	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only run head on one table")
		os.Exit(1)
		return
	}
	var reader *csv.Reader
	if len(moreArgs) == 1 {
		file, err := os.Open(moreArgs[0])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader = csv.NewReader(file)
	} else {
		reader = csv.NewReader(os.Stdin)
	}
	if strings.HasPrefix(numRowsStr, "+") {
		numRowsStr = strings.TrimPrefix(numRowsStr, "+")
		numRows, err := strconv.Atoi(numRowsStr)
		if err != nil {
			panic(err)
		}
		HeadFromBottom(reader, numRows)
	} else {
		numRows, err := strconv.Atoi(numRowsStr)
		if err != nil {
			panic(err)
		}
		HeadFromTop(reader, numRows)
	}
}
