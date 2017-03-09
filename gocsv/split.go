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

func Split(reader *csv.Reader, maxRows int, filenameBase string) {
	// Read and write header.
	header, err := reader.Read()
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
		row, err := reader.Read()
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
	moreArgs := fs.Args()
	if len(moreArgs) > 1 {
		fmt.Fprintln(os.Stderr, "Can only split one file")
		return
	}
	var reader *csv.Reader
	if len(moreArgs) == 1 {
		file, err := os.Open(moreArgs[0])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader = csv.NewReader(file)
		if filenameBase == "" {
			fileParts := strings.Split(moreArgs[0], ".")
			filenameBase = strings.Join(fileParts[:len(fileParts)-1], ".")
		}
	} else {
		reader = csv.NewReader(os.Stdin)
		if filenameBase == "" {
			filenameBase = "out"
		}
	}
	Split(reader, maxRows, filenameBase)
}
