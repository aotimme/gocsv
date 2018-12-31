package cmd

import (
	"flag"
	"fmt"
	"io"
	"strconv"
)

type DimensionsSubcommand struct {
	asCsv bool
}

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
	fs.BoolVar(&sub.asCsv, "csv", false, "Output results as CSV")
}

func (sub *DimensionsSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	GetDimensions(inputCsvs[0], sub.asCsv)
}

func GetDimensions(inputCsv *InputCsv, asCsv bool) {
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

	if asCsv {
		outputCsv := NewOutputCsvFromInputCsv(inputCsv)
		outputCsv.Write([]string{"Dimension", "Size"})
		outputCsv.Write([]string{"Rows", strconv.Itoa(numRows)})
		outputCsv.Write([]string{"Columns", strconv.Itoa(numColumns)})
	} else {
		fmt.Println("Dimensions:")
		fmt.Printf("  Rows: %d\n", numRows)
		fmt.Printf("  Columns: %d\n", numColumns)
	}
}
