package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
	"strconv"
)

type AutoincrementSubcommand struct {
	name    string
	seed    int
	prepend bool
}

func (sub *AutoincrementSubcommand) Name() string {
	return "autoincrement"
}
func (sub *AutoincrementSubcommand) Aliases() []string {
	return []string{"autoinc"}
}
func (sub *AutoincrementSubcommand) Description() string {
	return "Add a column of incrementing integers to a CSV."
}
func (sub *AutoincrementSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.name, "name", "ID", "Name of autoincrementing column")
	fs.IntVar(&sub.seed, "seed", 1, "Initial value of autoincrementing column")
	fs.BoolVar(&sub.prepend, "prepend", false, "Prepend the autoincrementing column (defaults to append)")
}

func (sub *AutoincrementSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	AutoIncrement(inputCsvs[0], sub.name, sub.seed, sub.prepend)
	err := inputCsvs[0].Close()
	if err != nil {
		panic(err)
	}
}

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
