package main

import (
	"flag"
	"fmt"
	"io"
)

type DimensionsSubcommand struct{}

func (sub *DimensionsSubcommand) Name() string {
	return "dimensions"
}
func (sub *DimensionsSubcommand) Aliases() []string {
	return []string{"dims"}
}
func (sub *DimensionsSubcommand) Description() string {
	return "Get the dimensions of a CSV."
}

func (sub *DimensionsSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	GetDimensions(inputCsvs[0])
}

func GetDimensions(inputCsv AbstractInputCsv) {
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	numColumns := len(header)

	numRows := 0
	for {
		_, err := inputCsv.Read()
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
