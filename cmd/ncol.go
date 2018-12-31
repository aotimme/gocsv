package cmd

import (
	"flag"
	"fmt"
	"io"
)

type NcolSubcommand struct {
}

func (sub *NcolSubcommand) Name() string {
	return "ncol"
}
func (sub *NcolSubcommand) Aliases() []string {
	return []string{}
}
func (sub *NcolSubcommand) Description() string {
	return "Get the number of columns in a CSV."
}
func (sub *NcolSubcommand) SetFlags(fs *flag.FlagSet) {
}

func (sub *NcolSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	GetNcol(inputCsvs[0])
}

func GetNcol(inputCsv *InputCsv) {
	// Be lenient when reading in the file.
	inputCsv.SetFieldsPerRecord(-1)

	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	numColumns := len(header)
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		if len(row) > numColumns {
			numColumns = len(row)
		}
	}
	fmt.Println(numColumns)
}
