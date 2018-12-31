package main

import (
	"flag"
	"io"
	"strconv"

	"github.com/alphagov/router/trie"
)

type UniqueSubcommand struct {
	columnsString string
	sorted        bool
	count         bool
}

func (sub *UniqueSubcommand) Name() string {
	return "unique"
}
func (sub *UniqueSubcommand) Aliases() []string {
	return []string{"uniq"}
}
func (sub *UniqueSubcommand) Description() string {
	return "Extract unique rows based upon certain columns."
}
func (sub *UniqueSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.columnsString, "columns", "", "Columns to use for comparison")
	fs.StringVar(&sub.columnsString, "c", "", "Columns to use for comparison (shorthand)")
	fs.BoolVar(&sub.sorted, "sorted", false, "Whether input CSV is already sorted")
	fs.BoolVar(&sub.count, "count", false, "Whether to append a Count column")
}

func (sub *UniqueSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsv(inputCsvs[0])
	sub.RunUnique(inputCsvs[0], outputCsv)
}

func (sub *UniqueSubcommand) RunUnique(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	var columns []string
	if sub.columnsString == "" {
		columns = make([]string, 0)
	} else {
		columns = GetArrayFromCsvString(sub.columnsString)
	}

	if sub.sorted {
		if sub.count {
			UniqueifySortedWithCount(inputCsv, outputCsvWriter, columns)
		} else {
			UniqueifySorted(inputCsv, outputCsvWriter, columns)
		}
	} else {
		if sub.count {
			UniqueifyUnsortedWithCount(inputCsv, outputCsvWriter, columns)
		} else {
			UniqueifyUnsorted(inputCsv, outputCsvWriter, columns)
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

func UniqueifySortedWithCount(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, columns []string) {
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}

	shellRow := make([]string, len(header)+1)

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	// Write header.
	copy(shellRow, header)
	shellRow[len(shellRow)-1] = "Count"
	outputCsvWriter.Write(shellRow)

	// Read and write first row.
	lastRow, err := inputCsv.Read()
	if err != nil {
		if err == io.EOF {
			return
		} else {
			ExitWithError(err)
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
				ExitWithError(err)
			}
		}
		if rowMatchesOnIndices(row, lastRow, columnIndices) {
			numInRun++
		} else {
			copy(shellRow, lastRow)
			shellRow[len(shellRow)-1] = strconv.Itoa(numInRun)
			outputCsvWriter.Write(shellRow)
			lastRow = row
			numInRun = 1
		}
	}
	copy(shellRow, lastRow)
	shellRow[len(shellRow)-1] = strconv.Itoa(numInRun)
	outputCsvWriter.Write(shellRow)
}

func UniqueifySorted(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, columns []string) {
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	// Write header.
	outputCsvWriter.Write(header)

	// Read and write first row.
	lastRow, err := inputCsv.Read()
	if err != nil {
		if err == io.EOF {
			return
		} else {
			ExitWithError(err)
		}
	}
	outputCsvWriter.Write(lastRow)

	// Write unique rows in order.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		if !rowMatchesOnIndices(row, lastRow, columnIndices) {
			lastRow = row
			outputCsvWriter.Write(row)
		}
	}
}

func UniqueifyUnsorted(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, columns []string) {
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	// Write header.
	outputCsvWriter.Write(header)

	seenRowsTrie := trie.NewTrie()
	lastRowArray := make([]string, len(columnIndices))

	// Write unique rows in order.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		for i, columnIndex := range columnIndices {
			lastRowArray[i] = row[columnIndex]
		}
		_, ok := seenRowsTrie.Get(lastRowArray)
		if !ok {
			seenRowsTrie.Set(lastRowArray, true)
			outputCsvWriter.Write(row)
		}
	}
}

func UniqueifyUnsortedWithCount(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, columns []string) {
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

	// Write header.
	outputCsvWriter.Write(shellRow)

	// Write unique rows with count.
	for rowIndex, row := range imc.rows {
		count, ok := rowIndexToCount[rowIndex]
		if ok {
			copy(shellRow, row)
			shellRow[len(shellRow)-1] = strconv.Itoa(count)
			outputCsvWriter.Write(shellRow)
		}
	}
}
