package main

import (
	"encoding/csv"
	"errors"
	"os"
)

type AbstractInputCsv interface {
	Close() error
	Read() ([]string, error)
	ReadAll() ([][]string, error)
	Reader() *csv.Reader
	Name() string
	Filename() string
}

type InputCsv struct {
	reader *csv.Reader
}

func (ic *InputCsv) Reader() *csv.Reader {
	return ic.reader
}

func (ic *InputCsv) Read() ([]string, error) {
	return ic.reader.Read()
}

func (ic *InputCsv) ReadAll() ([][]string, error) {
	return ic.reader.ReadAll()
}

type StandardInputCsv struct {
	InputCsv
}

func (sic *StandardInputCsv) Close() error {
	return nil
}

func (sic *StandardInputCsv) Name() string {
	return "-"
}

func (sic *StandardInputCsv) Filename() string {
	return "-"
}

func NewStandardInputCsv() (sic *StandardInputCsv) {
	sic = new(StandardInputCsv)
	sic.reader = csv.NewReader(os.Stdin)
	return
}

type FileInputCsv struct {
	InputCsv
	filename string
	file     *os.File
}

func NewFileInputCsv(filename string) (fic *FileInputCsv, err error) {
	fic = new(FileInputCsv)
	fic.filename = filename
	fic.file, err = os.Open(filename)
	if err != nil {
		return
	}
	fic.reader = csv.NewReader(fic.file)
	return
}

func (fic *FileInputCsv) Close() error {
	return fic.file.Close()
}

func (fic *FileInputCsv) Name() string {
	return GetBaseFilenameWithoutExtension(fic.filename)
}

func (fic *FileInputCsv) Filename() string {
	return fic.filename
}

func GetInputCsvs(filenames []string, numInputCsvs int) (csvs []AbstractInputCsv, err error) {
	hasDash := false
	for _, filename := range filenames {
		if filename == "-" {
			hasDash = true
			break
		}
	}
	if numInputCsvs == -1 {
		if len(filenames) == 0 {
			csvs = make([]AbstractInputCsv, 1)
			csvs[0] = NewStandardInputCsv()
			return
		} else {
			csvs = make([]AbstractInputCsv, len(filenames))
			for i, filename := range filenames {
				csvs[i], err = NewAbstractInputCsv(filename)
				if err != nil {
					return
				}
			}
			return
		}
	} else {
		csvs = make([]AbstractInputCsv, numInputCsvs)
		if len(filenames) > numInputCsvs {
			err = errors.New("Too many files for command")
			return
		}
		if len(filenames) == numInputCsvs {
			for i, filename := range filenames {
				csvs[i], err = NewAbstractInputCsv(filename)
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
			csvs[0] = NewStandardInputCsv()
			for i, filename := range filenames {
				csvs[i+1], err = NewAbstractInputCsv(filename)
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

func NewAbstractInputCsv(filename string) (aic AbstractInputCsv, err error) {
	if filename == "-" {
		aic = NewStandardInputCsv()
	} else {
		aic, err = NewFileInputCsv(filename)
	}
	return
}
