package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

func ChangeDelimiter(reader *csv.Reader, inputDelimiter, outputDelimiter string) {
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

func RunDelimiter(args []string) {
	fs := flag.NewFlagSet("delimiter", flag.ExitOnError)
	var inputDelimiter, outputDelimiter string
	fs.StringVar(&inputDelimiter, "input", "", "Input delimiter")
	fs.StringVar(&inputDelimiter, "i", "", "Input delimiter (shorthand)")
	fs.StringVar(&outputDelimiter, "output", "", "Output delimiter")
	fs.StringVar(&outputDelimiter, "o", "", "Output delimiter (shorthand)")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only change delimiters on one table")
		os.Exit(1)
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
	ChangeDelimiter(reader, inputDelimiter, outputDelimiter)
}
