package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
	"strconv"
)

func AutoIncrement(inputCsv AbstractInputCsv, name string, seed int, prepend bool) {
	writer := csv.NewWriter(os.Stdout)

	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	numInputColumns := len(header)
	shellRow := make([]string, numInputColumns+1)
	if prepend {
		shellRow[0] = name
		for i, elem := range header {
			shellRow[i+1] = elem
		}
	} else {
		copy(shellRow, header)
		shellRow[numInputColumns] = name
	}
	writer.Write(shellRow)
	writer.Flush()

	// Write rows with autoincrement.
	inc := seed
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		incStr := strconv.Itoa(inc)
		if prepend {
			shellRow[0] = incStr
			for i, elem := range row {
				shellRow[i+1] = elem
			}
		} else {
			copy(shellRow, row)
			shellRow[numInputColumns] = incStr
		}
		inc++
		writer.Write(shellRow)
		writer.Flush()
	}
}

func RunAutoIncrement(args []string) {
	fs := flag.NewFlagSet("autoincrement", flag.ExitOnError)
	var name string
	var seed int
	var prepend bool
	fs.StringVar(&name, "name", "ID", "Name of autoincrementing column")
	fs.IntVar(&seed, "seed", 1, "Initial value of autoincrementing column")
	fs.BoolVar(&prepend, "prepend", false, "Prepend the autoincrementing column (defaults to append)")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}
	AutoIncrement(inputCsvs[0], name, seed, prepend)
	err = inputCsvs[0].Close()
	if err != nil {
		panic(err)
	}
}
