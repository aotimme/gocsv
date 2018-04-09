package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func Split(inputCsv AbstractInputCsv, maxRows int, filenameBase string) {
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
		panic(err)
	}

	fileNumber := 1
	numRowsWritten := 0
	curFilename := filenameBase + "-" + strconv.Itoa(fileNumber) + ".csv"
	curFile, err := os.Create(curFilename)
	if err != nil {
		panic(err)
	}
	defer curFile.Close()
	writer := csv.NewWriter(curFile)
	writer.Write(header)
	writer.Flush()

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		// Switch to the next file.
		if numRowsWritten == maxRows {
			fileNumber++
			numRowsWritten = 0
			curFilename = filenameBase + "-" + strconv.Itoa(fileNumber) + ".csv"
			curFile, err = os.Create(curFilename)
			if err != nil {
				panic(err)
			}
			defer curFile.Close()
			writer = csv.NewWriter(curFile)
			writer.Write(header)
			writer.Flush()
		}

		writer.Write(row)
		writer.Flush()
		numRowsWritten++
	}
}

func RunSplit(args []string) {
	fs := flag.NewFlagSet("split", flag.ExitOnError)
	var maxRows int
	var filenameBase string
	fs.IntVar(&maxRows, "max-rows", 0, "Maximum number of rows per CSV.")
	fs.StringVar(&filenameBase, "filename-base", "", "Base of filenames for output.")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	if maxRows < 1 {
		fmt.Fprintln(os.Stderr, "Invalid parameter for --max-rows")
		os.Exit(1)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}
	Split(inputCsvs[0], maxRows, filenameBase)
}
