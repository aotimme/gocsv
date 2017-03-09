package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func DescribeCsv(reader *csv.Reader) {
	imc := NewInMemoryCsv(reader)

	numRows := imc.NumRows()
	numColumns := imc.NumColumns()

	fmt.Println("Dimensions:")
	fmt.Printf("  Rows: %d\n", numRows)
	fmt.Printf("  Columns: %d\n", numColumns)
	fmt.Println("Columns:")

	for i := 0; i < numColumns; i++ {
		columnType := imc.InferType(i)
		fmt.Printf("  %d: %s\n", i+1, imc.header[i])
		fmt.Printf("    Type: %s\n", ColumnTypeToString(columnType))
	}
}

func RunDescribe(args []string) {
	fs := flag.NewFlagSet("describe", flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only describe one CSV")
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

	DescribeCsv(reader)
}
