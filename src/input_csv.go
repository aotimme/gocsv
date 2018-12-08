package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"gocsv/csv"
)

type InputCsv struct {
	file      *os.File
	filename  string
	reader    *csv.Reader
	bufReader *bufio.Reader
	hasBom    bool
}

func NewInputCsv(filename string) (ic *InputCsv, err error) {
	ic = new(InputCsv)
	ic.filename = filename
	if filename == "-" {
		ic.file = os.Stdin
	} else {
		ic.file, err = os.Open(filename)
		if err != nil {
			return
		}
	}
	ic.bufReader = bufio.NewReader(ic.file)
	ic.reader = csv.NewReader(ic.bufReader)
	err = ic.handleBom()
	return
}

func (ic *InputCsv) handleBom() error {
	bomRune, _, err := ic.bufReader.ReadRune()
	if err != nil && err != io.EOF {
		return err
	}
	if err != io.EOF && bomRune == BOM_RUNE {
		ic.hasBom = true
	} else {
		ic.bufReader.UnreadRune()
	}
	return nil
}

func (ic *InputCsv) Close() error {
	return ic.file.Close()
}

func (ic *InputCsv) SetFieldsPerRecord(fieldsPerRecord int) {
	ic.reader.FieldsPerRecord = fieldsPerRecord
}

func (ic *InputCsv) SetLazyQuotes(lazyQuotes bool) {
	ic.reader.LazyQuotes = lazyQuotes
}

func (ic *InputCsv) SetDelimiter(delimiter rune) {
	ic.reader.Comma = delimiter
}

func (ic *InputCsv) Reader() *csv.Reader {
	return ic.reader
}

func (ic *InputCsv) Read() (row []string, err error) {
	return ic.reader.Read()
}

func (ic *InputCsv) ReadAll() (rows [][]string, err error) {
	return ic.reader.ReadAll()
}

func (ic *InputCsv) Name() string {
	if ic.filename == "-" {
		return "stdin"
	} else {
		return GetBaseFilenameWithoutExtension(ic.filename)
	}
}

func (ic *InputCsv) Filename() string {
	return ic.filename
}

func GetInputCsvsOrPanic(filenames []string, numInputCsvs int) (csvs []*InputCsv) {
	csvs, err := GetInputCsvs(filenames, numInputCsvs)
	if err != nil {
		ExitWithError(err)
	}
	return
}

func GetInputCsvs(filenames []string, numInputCsvs int) (csvs []*InputCsv, err error) {
	hasDash := false
	for _, filename := range filenames {
		if filename == "-" {
			hasDash = true
			break
		}
	}
	if numInputCsvs == -1 {
		if len(filenames) == 0 {
			csvs = make([]*InputCsv, 1)
			csvs[0], err = NewInputCsv("-")
			return
		} else {
			csvs = make([]*InputCsv, len(filenames))
			for i, filename := range filenames {
				csvs[i], err = NewInputCsv(filename)
				if err != nil {
					return
				}
			}
			return
		}
	} else {
		csvs = make([]*InputCsv, numInputCsvs)
		if len(filenames) > numInputCsvs {
			err = errors.New("Too many files for command")
			return
		}
		if len(filenames) == numInputCsvs {
			for i, filename := range filenames {
				csvs[i], err = NewInputCsv(filename)
				if err != nil {
					return
				}
			}
			return
		}
		if len(filenames) == numInputCsvs-1 {
			if hasDash {
				err = errors.New("Too few inputs specified")
				return
			}
			csvs[0], err = NewInputCsv("-")
			if err != nil {
				return
			}
			for i, filename := range filenames {
				csvs[i+1], err = NewInputCsv(filename)
				if err != nil {
					return
				}
			}
			return
		}
		err = errors.New("Too few inputs specified")
		return
	}
}
