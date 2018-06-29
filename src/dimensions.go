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
func (sub *DimensionsSubcommand) SetFlags(fs *flag.FlagSet) {
}

func (sub *DimensionsSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	GetDimensions(inputCsvs[0])
}

func GetDimensions(inputCsv *InputCsv) {
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	numColumns := len(header)

	numRows := 0
	for {
		_, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		numRows++
	}

	fmt.Println("Dimensions:")
	fmt.Printf("  Rows: %d\n", numRows)
	fmt.Printf("  Columns: %d\n", numColumns)
}
