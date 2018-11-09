package main

import (
	"flag"
	"io"
)

type ZipSubcommand struct {
}

func (sub *ZipSubcommand) Name() string {
	return "zip"
}
func (sub *ZipSubcommand) Aliases() []string {
	return []string{}
}
func (sub *ZipSubcommand) Description() string {
	return "Zip multiple CSVs into one CSV."
}
func (sub *ZipSubcommand) SetFlags(fs *flag.FlagSet) {
}

func (sub *ZipSubcommand) Run(args []string) {
	filenames := args
	inputCsvs := GetInputCsvsOrPanic(filenames, -1)
	ZipFiles(inputCsvs)
}

func ZipFiles(inputCsvs []*InputCsv) {
	outputCsv := NewOutputCsvFromInputCsvs(inputCsvs)

	numCsvs := len(inputCsvs)

	numColumns := 0
	offsets := make([]int, numCsvs+1)
	offsets[0] = 0
	headers := make([][]string, numCsvs)
	for i, inputCsv := range inputCsvs {
		header, err := inputCsv.Read()
		if err != nil {
			ExitWithError(err)
		}
		headers[i] = header
		numColumns += len(header)
		offsets[i+1] = numColumns
	}

	shellRow := make([]string, numColumns)

	for i, header := range headers {
		start := offsets[i]
		end := offsets[i+1]
		copy(shellRow[start:end], header)
	}
	outputCsv.Write(shellRow)

	isInputCsvComplete := make([]bool, numCsvs)
	numCsvsComplete := 0
	// Go through the files
	for true {
		for i, inputCsv := range inputCsvs {
			if isInputCsvComplete[i] {
				continue
			}
			start := offsets[i]
			end := offsets[i+1]
			row, err := inputCsv.Read()
			if err != nil {
				if err == io.EOF {
					isInputCsvComplete[i] = true
					numCsvsComplete++
					copy(shellRow[start:end], make([]string, end-start))
				} else {
					ExitWithError(err)
				}
			} else {
				copy(shellRow[start:end], row)
			}
			if numCsvsComplete == numCsvs {
				break
			}
		}
		if numCsvsComplete == numCsvs {
			break
		}
		outputCsv.Write(shellRow)
	}
}
