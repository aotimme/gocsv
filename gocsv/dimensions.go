package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

func GetDimensions(reader *csv.Reader) {
	header, err := reader.Read()
	if err != nil {
		panic(err)
	}
	numColumns := len(header)

	numRows := 0
	for {
		_, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		numRows++
	}

	fmt.Println("Dimensions:")
	fmt.Printf("  Rows: %d\n", numRows)
	fmt.Printf("  Columns: %d\n", numColumns)
}

func RunDimensions(args []string) {
	fs := flag.NewFlagSet("dimensions", flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only get dimensions on one CSV")
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

	GetDimensions(reader)
}
