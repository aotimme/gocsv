package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type SplitSubcommand struct {
	maxRows      int
	filenameBase string
}

func (sub *SplitSubcommand) Name() string {
	return "split"
}
func (sub *SplitSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SplitSubcommand) Description() string {
	return "Split a CSV into multiple files."
}
func (sub *SplitSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.IntVar(&sub.maxRows, "max-rows", 0, "Maximum number of rows per CSV.")
	fs.StringVar(&sub.filenameBase, "filename-base", "", "Base of filenames for output.")
}

func (sub *SplitSubcommand) Run(args []string) {
	if sub.maxRows < 1 {
		fmt.Fprintln(os.Stderr, "Invalid parameter for --max-rows")
		os.Exit(1)
	}

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	Split(inputCsvs[0], sub.maxRows, sub.filenameBase)
}

func Split(inputCsv *InputCsv, maxRows int, filenameBase string) {
	if filenameBase == "" {
		inputFilename := inputCsv.Filename()
		if inputFilename == "-" {
			filenameBase = "out"
		} else {
			fileParts := strings.Split(inputFilename, ".")
			filenameBase = strings.Join(fileParts[:len(fileParts)-1], ".")
		}
	}

	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}

	fileNumber := 1
	numRowsWritten := 0
	curFilename := filenameBase + "-" + strconv.Itoa(fileNumber) + ".csv"
	curFile, err := os.Create(curFilename)
	if err != nil {
		ExitWithError(err)
	}
	defer curFile.Close()

	outputCsv := NewFileOutputCsvFromInputCsv(inputCsv, curFile)
	outputCsv.Write(header)

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		// Switch to the next file.
		if numRowsWritten == maxRows {
			fileNumber++
			numRowsWritten = 0
			curFilename = filenameBase + "-" + strconv.Itoa(fileNumber) + ".csv"
			curFile, err = os.Create(curFilename)
			if err != nil {
				ExitWithError(err)
			}
			defer curFile.Close()
			outputCsv = NewFileOutputCsvFromInputCsv(inputCsv, curFile)
			outputCsv.Write(header)
		}

		outputCsv.Write(row)
		numRowsWritten++
	}
}
