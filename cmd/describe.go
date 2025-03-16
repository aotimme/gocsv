package cmd

import (
	"flag"
	"fmt"
)

type DescribeSubcommand struct{}

func (sub *DescribeSubcommand) Name() string {
	return "describe"
}
func (sub *DescribeSubcommand) Aliases() []string {
	return []string{}
}
func (sub *DescribeSubcommand) Description() string {
	return "Get basic information about a CSV."
}
func (sub *DescribeSubcommand) SetFlags(fs *flag.FlagSet) {
}

func (sub *DescribeSubcommand) Run(args []string) {
	useTimeLayoutEnvVar()

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	DescribeCsv(inputCsvs[0])
}

func DescribeCsv(inputCsv *InputCsv) {
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
