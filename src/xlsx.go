package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"./csv"

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

func ConvertXlsxSheetToDirectory(dirname string, sheet *xlsx.Sheet) {
	filename := fmt.Sprintf("%s/%s.csv", dirname, sheet.Name)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	for _, row := range sheet.Rows {
		csvRow := make([]string, 0)
		for _, cell := range row.Cells {
			cellValue, err := cell.FormattedValue()
			if err != nil {
				panic(err)
			}
			csvRow = append(csvRow, cellValue)
		}
		writer.Write(csvRow)
		writer.Flush()
	}
}

func ConvertXlsxFull(filename, dirname string) {
	xlsxFile, err := xlsx.OpenFile(filename)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(dirname, os.ModeDir|0755)
	for _, sheet := range xlsxFile.Sheets {
		ConvertXlsxSheetToDirectory(dirname, sheet)
	}
}

func ConvertXlsxSheet(filename, sheetName string) {
	xlsxFile, err := xlsx.OpenFile(filename)
	if err != nil {
		panic(err)
	}

	sheetNames := make([]string, len(xlsxFile.Sheets))
	for i, sheet := range xlsxFile.Sheets {
		sheetNames[i] = sheet.Name
	}
	sheetIndex := GetIndexForColumn(sheetNames, sheetName)
	if sheetIndex == -1 {
		panic(errors.New("Could not find sheet from sheet name"))
	}

	sheet := xlsxFile.Sheets[sheetIndex]
	writer := csv.NewWriter(os.Stdout)
	for _, row := range sheet.Rows {
		csvRow := make([]string, 0)
		for _, cell := range row.Cells {
			cellValue, err := cell.FormattedValue()
			if err != nil {
				panic(err)
			}
			csvRow = append(csvRow, cellValue)
		}
		writer.Write(csvRow)
		writer.Flush()
	}
}

func ListXlxsSheets(filename string) {
	xlsxFile, err := xlsx.OpenFile(filename)
	if err != nil {
		panic(err)
	}

	for i, sheet := range xlsxFile.Sheets {
		fmt.Printf("%d: %s\n", i+1, sheet.Name)
	}
}
