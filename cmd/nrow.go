package cmd

import (
	"flag"
	"fmt"
	"io"
)

type NrowSubcommand struct {
}

func (sub *NrowSubcommand) Name() string {
	return "nrow"
}
func (sub *NrowSubcommand) Aliases() []string {
	return []string{}
}
func (sub *NrowSubcommand) Description() string {
	return "Get the number of rows in a CSV."
}
func (sub *NrowSubcommand) SetFlags(fs *flag.FlagSet) {
}

func (sub *NrowSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	GetNrow(inputCsvs[0])
}

func GetNrow(inputCsv *InputCsv) {
	_, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}

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
	fmt.Println(numRows)
}
