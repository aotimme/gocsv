package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

type XlsxSubcommand struct {
	listSheets bool
	dirname    string
	sheet      string
}

func (sub *XlsxSubcommand) Name() string {
	return "xlsx"
}
func (sub *XlsxSubcommand) Aliases() []string {
	return []string{}
}
func (sub *XlsxSubcommand) Description() string {
	return "Convert sheets of a XLSX file to CSV."
}
func (sub *XlsxSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&sub.listSheets, "list-sheets", false, "List sheets in file")
	fs.StringVar(&sub.dirname, "dirname", "", "Name of folder to output sheets to")
	fs.StringVar(&sub.sheet, "sheet", "", "Name of sheet to convert")
}

func (sub *XlsxSubcommand) Run(args []string) {
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Can only convert one file")
		os.Exit(1)
	} else if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Cannot convert file from stdin")
		os.Exit(1)
	}
	filename := args[0]
	if sub.listSheets {
		ListXlxsSheets(filename)
	} else {
		if sub.sheet == "" {
			if sub.dirname == "" {
				fileParts := strings.Split(filename, ".")
				sub.dirname = strings.Join(fileParts[:len(fileParts)-1], ".")
			}
			ConvertXlsxFull(filename, sub.dirname)
		} else {
			ConvertXlsxSheet(filename, sub.sheet)
		}
	}
}

func ConvertXlsxFull(filename, dirname string) {
	xlsxFile, err := xlsx.OpenFile(filename)
	if err != nil {
		ExitWithError(err)
	}
	err = os.Mkdir(dirname, os.ModeDir|0755)
	if err != nil {
		ExitWithError(err)
	}
	for _, sheet := range xlsxFile.Sheets {
		ConvertXlsxSheetToDirectory(dirname, sheet)
	}
}

func ConvertXlsxSheetToDirectory(dirname string, sheet *xlsx.Sheet) {
	filename := fmt.Sprintf("%s/%s.csv", dirname, sheet.Name)

	file, err := os.Create(filename)
	if err != nil {
		ExitWithError(err)
	}
	defer file.Close()
	outputCsv := NewOutputCsvFromFile(file)
	WriteSheetToOutputCsv(sheet, outputCsv)
}

func ConvertXlsxSheet(filename, sheetName string) {
	xlsxFile, err := xlsx.OpenFile(filename)
	if err != nil {
		ExitWithError(err)
	}

	sheetNames := make([]string, len(xlsxFile.Sheets))
	for i, sheet := range xlsxFile.Sheets {
		sheetNames[i] = sheet.Name
	}
	sheetIndex := GetIndexForColumn(sheetNames, sheetName)
	if sheetIndex == -1 {
		ExitWithError(errors.New("Could not find sheet from sheet name"))
	}

	sheet := xlsxFile.Sheets[sheetIndex]
	outputCsv := NewOutputCsv()
	WriteSheetToOutputCsv(sheet, outputCsv)
}

func WriteSheetToOutputCsv(sheet *xlsx.Sheet, outputCsv *OutputCsv) {
	for _, row := range sheet.Rows {
		csvRow := make([]string, 0)
		for _, cell := range row.Cells {
			// We only care about the string value of the cell,
			// so just ignore any error.
			cellValue, _ := cell.FormattedValue()
			csvRow = append(csvRow, cellValue)
		}
		outputCsv.Write(csvRow)
	}
}

func ListXlxsSheets(filename string) {
	xlsxFile, err := xlsx.OpenFile(filename)
	if err != nil {
		ExitWithError(err)
	}

	for i, sheet := range xlsxFile.Sheets {
		fmt.Printf("%d: %s\n", i+1, sheet.Name)
	}
}
