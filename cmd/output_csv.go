package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type OutputCsvWriter interface {
	Write(row []string) error
}

type OutputCsv struct {
	writeBom         bool
	hasWrittenHeader bool
	csvWriter        *csv.Writer
	file             *os.File
	writeRaw         bool
	flushEveryRow    bool
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
	oc.file = file
	oc.csvWriter = csv.NewWriter(file)
	if value := os.Getenv("GOCSV_FLUSH_EVERY_ROW"); value != "" {
		var err error
		oc.flushEveryRow, err = strconv.ParseBool(value)
		if err != nil {
			ExitWithError(fmt.Errorf("invalid GOCSV_FLUSH_EVERY_ROW value %q: %w", value, err))
		}
	}
	delimiter := os.Getenv("GOCSV_DELIMITER")
	if delimiter != "" {
		oc.csvWriter.Comma = GetDelimiterFromStringOrPanic(delimiter)
	}
	return
}

// Flush writes any buffered CSV data to the output file and returns any write error.
// Call it when finished writing, before closing the file.
func (oc *OutputCsv) Flush() error {
	oc.csvWriter.Flush()
	return oc.csvWriter.Error()
}

// FlushOrPanic flushes buffered CSV data and reports any error via ExitWithError.
func (oc *OutputCsv) FlushOrPanic() {
	if err := oc.Flush(); err != nil {
		ExitWithError(err)
	}
}

func (oc *OutputCsv) SetDelimiter(delimiter rune) {
	oc.csvWriter.Comma = delimiter
}

func (oc *OutputCsv) SetWriteRaw(writeRaw bool) {
	oc.writeRaw = writeRaw
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
	if oc.writeRaw && len(row) == 1 {
		oc.file.WriteString(row[0] + "\n")
		return
	}
	err = oc.csvWriter.Write(row)
	if err != nil {
		return
	}
	if oc.flushEveryRow {
		return oc.Flush()
	}
	return
}
