package main

import (
	"flag"
	"fmt"
)

func DescribeCsv(inputCsv AbstractInputCsv) {
	imc := NewInMemoryCsvFromInputCsv(inputCsv)

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

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}
	DescribeCsv(inputCsvs[0])
}
