package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
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
		if sub.sheet != "" && sub.dirname != "" {
			fmt.Fprintln(os.Stderr, "Cannot use --sheet and --dirname together")
			os.Exit(1)
		}
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
	f, err := excelize.OpenFile(filename)
	if err != nil {
		ExitWithError(err)
	}
	err = os.Mkdir(dirname, os.ModeDir|0755)
	if err != nil {
		ExitWithError(err)
	}
	for _, sheetName := range f.GetSheetList() {
		ConvertXlsxSheetToDirectory(f, dirname, sheetName)
	}
}

func ConvertXlsxSheetToDirectory(f *excelize.File, dirname string, sheetName string) {
	filename := fmt.Sprintf("%s/%s.csv", dirname, sheetName)

	file, err := os.Create(filename)
	if err != nil {
		ExitWithError(err)
	}
	defer file.Close()
	outputCsv := NewOutputCsvFromFile(file)
	writeRowsToOutputCsv(outputCsv, f, sheetName)
}

func ConvertXlsxSheet(filename, sheetName string) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		ExitWithError(err)
	}
	sheetNames := f.GetSheetList()
	// Use `GetIndexForColumn` so the sheet can be specified
	// by name or by index.
	sheetIndex := GetIndexForColumn(sheetNames, sheetName)
	if sheetIndex == -1 {
		if sheetIndex == -1 {
			ExitWithError(errors.New("could not find sheet from sheet name"))
		}
	}
	trueSheetName := sheetNames[sheetIndex]
	outputCsv := NewOutputCsv()
	writeRowsToOutputCsv(outputCsv, f, trueSheetName)
}

func writeRowsToOutputCsv(outputCsv *OutputCsv, f *excelize.File, sheetName string) {
	rows, err := f.Rows(sheetName)
	if err != nil {
		ExitWithError(err)
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			ExitWithError(err)
		}
		outputCsv.Write(row)
	}
	if err = rows.Close(); err != nil {
		ExitWithError(err)
	}
}

func ListXlxsSheets(filename string) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		ExitWithError(err)
	}

	for i, sheetName := range f.GetSheetList() {
		fmt.Printf("%d: %s\n", i+1, sheetName)
	}
}
