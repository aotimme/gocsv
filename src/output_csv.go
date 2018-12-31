package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type OutputCsvWriter interface {
	Write(row []string) error
}

type OutputCsv struct {
	writeBom         bool
	hasWrittenHeader bool
	writer           *csv.Writer
}

func NewOutputCsvFromInputCsv(inputCsv *InputCsv) (oc *OutputCsv) {
	return NewOutputCsvFromInputCsvs([]*InputCsv{inputCsv})
}

func NewOutputCsvFromInputCsvs(inputCsvs []*InputCsv) (oc *OutputCsv) {
	return NewOutputCsvFromInputCsvsAndFile(inputCsvs, os.Stdout)
}

func NewFileOutputCsvFromInputCsv(inputCsv *InputCsv, file *os.File) (oc *OutputCsv) {
	return NewOutputCsvFromInputCsvsAndFile([]*InputCsv{inputCsv}, file)
}

func NewOutputCsvFromInputCsvsAndFile(inputCsvs []*InputCsv, file *os.File) (oc *OutputCsv) {
	oc = NewOutputCsvFromFile(file)
	// If _any_ of the input CSVs has a BOM, then conserve the BOM.
	for _, inputCsv := range inputCsvs {
		if inputCsv.hasBom {
			oc.writeBom = true
			break
		}
	}
	return
}

func NewOutputCsv() (oc *OutputCsv) {
	return NewOutputCsvFromFile(os.Stdout)
}

func NewOutputCsvFromFile(file *os.File) (oc *OutputCsv) {
	oc = new(OutputCsv)
	oc.writer = csv.NewWriter(file)
	return
}

func (oc *OutputCsv) SetDelimiter(delimiter rune) {
	oc.writer.Comma = delimiter
}

func (oc *OutputCsv) Write(row []string) error {
	if !oc.hasWrittenHeader {
		oc.hasWrittenHeader = true
		if oc.writeBom {
			rowCopy := make([]string, len(row))
			copy(rowCopy, row)
			rowCopy[0] = fmt.Sprintf("%s%s", BOM_STRING, row[0])
			return oc.writeRow(rowCopy)
		}
	}
	return oc.writeRow(row)
}

func (oc *OutputCsv) writeRow(row []string) (err error) {
	err = oc.writer.Write(row)
	if err != nil {
		return
	}
	// It is less efficient to flush after every write, but doing so
	// keeps the output flowing. Otherwise it could look "jumpy" or
	// like it's not working at times when there is no visible output
	// while working on a large file.
	oc.writer.Flush()
	return
}
