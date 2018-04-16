package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
	"strconv"

	"github.com/alphagov/router/trie"
)

type UniqueSubcommand struct{}

func (sub *UniqueSubcommand) Name() string {
	return "unique"
}
func (sub *UniqueSubcommand) Aliases() []string {
	return []string{"uniq"}
}
func (sub *UniqueSubcommand) Description() string {
	return "Extract unique rows based upon certain columns."
}

func (sub *UniqueSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	var columnsString string
	var sorted, count bool
	fs.StringVar(&columnsString, "columns", "", "Columns to use for comparison")
	fs.StringVar(&columnsString, "c", "", "Columns to use for comparison (shorthand)")
	fs.BoolVar(&sorted, "sorted", false, "Whether input CSV is already sorted")
	fs.BoolVar(&count, "count", false, "Whether to append a Count column")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	var columns []string
	if columnsString == "" {
		columns = make([]string, 0)
	} else {
		columns = GetArrayFromCsvString(columnsString)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	if sorted {
		if count {
			UniqueifySortedWithCount(inputCsvs[0], columns)
		} else {
			UniqueifySorted(inputCsvs[0], columns)
		}
	} else {
		if count {
			UniqueifyUnsortedWithCount(inputCsvs[0], columns)
		} else {
			UniqueifyUnsorted(inputCsvs[0], columns)
		}
	}
}

func rowMatchesOnIndices(rowA, rowB []string, columnIndices []int) bool {
	for _, columnIndex := range columnIndices {
		if rowA[columnIndex] != rowB[columnIndex] {
			return false
		}
	}
	return true
}

func UniqueifySortedWithCount(inputCsv AbstractInputCsv, columns []string) {
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}

	shellRow := make([]string, len(header)+1)

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	copy(shellRow, header)
	shellRow[len(shellRow)-1] = "Count"
	writer.Write(shellRow)
	writer.Flush()

	// Read and write first row.
	lastRow, err := inputCsv.Read()
	if err != nil {
		if err == io.EOF {
			return
		} else {
			panic(err)
		}
	}
	numInRun := 1

	// Write unique rows in order.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		if rowMatchesOnIndices(row, lastRow, columnIndices) {
			numInRun++
		} else {
			copy(shellRow, lastRow)
			shellRow[len(shellRow)-1] = strconv.Itoa(numInRun)
			writer.Write(shellRow)
			writer.Flush()
			lastRow = row
			numInRun = 1
		}
	}
	copy(shellRow, lastRow)
	shellRow[len(shellRow)-1] = strconv.Itoa(numInRun)
	writer.Write(shellRow)
	writer.Flush()
}

func UniqueifySorted(inputCsv AbstractInputCsv, columns []string) {
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	writer.Write(header)
	writer.Flush()

	// Read and write first row.
	lastRow, err := inputCsv.Read()
	if err != nil {
		if err == io.EOF {
			return
		} else {
			panic(err)
		}
	}
	writer.Write(lastRow)
	writer.Flush()

	// Write unique rows in order.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		if !rowMatchesOnIndices(row, lastRow, columnIndices) {
			lastRow = row
			writer.Write(row)
			writer.Flush()
		}
	}
}

func UniqueifyUnsorted(inputCsv AbstractInputCsv, columns []string) {
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	writer.Write(header)
	writer.Flush()

	seenRowsTrie := trie.NewTrie()
	lastRowArray := make([]string, len(columnIndices))

	// Write unique rows in order.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		for i, columnIndex := range columnIndices {
			lastRowArray[i] = row[columnIndex]
		}
		_, ok := seenRowsTrie.Get(lastRowArray)
		if !ok {
			seenRowsTrie.Set(lastRowArray, true)
			writer.Write(row)
			writer.Flush()
		}
	}
}

func UniqueifyUnsortedWithCount(inputCsv AbstractInputCsv, columns []string) {
	imc := NewInMemoryCsvFromInputCsv(inputCsv)

	columnIndices := GetIndicesForColumnsOrPanic(imc.header, columns)

	rowIndexToCount := make(map[int]int)
	seenRowsTrie := trie.NewTrie()

	lastRowArray := make([]string, len(columnIndices))

	for rowIndex, row := range imc.rows {
		for i, columnIndex := range columnIndices {
			lastRowArray[i] = row[columnIndex]
		}
		val, ok := seenRowsTrie.Get(lastRowArray)
		if ok {
			previousRowIndex := val.(int)
			rowIndexToCount[previousRowIndex] = rowIndexToCount[previousRowIndex] + 1
		} else {
			previousRowIndex := rowIndex
			seenRowsTrie.Set(lastRowArray, previousRowIndex)
			rowIndexToCount[previousRowIndex] = 1
		}
	}

	shellRow := make([]string, len(imc.header)+1)
	copy(shellRow, imc.header)
	shellRow[len(shellRow)-1] = "Count"

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	writer.Write(shellRow)
	writer.Flush()

	// Write unique rows with count.
	for rowIndex, row := range imc.rows {
		count, ok := rowIndexToCount[rowIndex]
		if ok {
			copy(shellRow, row)
			shellRow[len(shellRow)-1] = strconv.Itoa(count)
			writer.Write(shellRow)
			writer.Flush()
		}
	}
}
