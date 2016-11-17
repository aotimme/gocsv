package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "os"
  "strings"
  "github.com/tealeg/xlsx"
)

func ConvertXLSXSheetToDirectory(dirname string, sheet *xlsx.Sheet) {
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

func ConvertXLSXFull(filename, dirname string) {
  xlsxFile, err := xlsx.OpenFile(filename)
  if err != nil {
    panic(err)
  }
  err = os.Mkdir(dirname, os.ModeDir | 0755)
  for _, sheet := range xlsxFile.Sheets {
    ConvertXLSXSheetToDirectory(dirname, sheet)
  }
}

func ConvertXLSXSheet(filename, sheetName string) {
  xlsxFile, err := xlsx.OpenFile(filename)
  if err != nil {
    panic(err)
  }

  sheetNames := make([]string, len(xlsxFile.Sheets))
  for i, sheet := range xlsxFile.Sheets {
    sheetNames[i] = sheet.Name
  }
  sheetIndex := GetColumnIndexOrPanic(sheetNames, sheetName)

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

func ListXLXSSheets(filename string) {
  xlsxFile, err := xlsx.OpenFile(filename)
  if err != nil {
    panic(err)
  }

  for i, sheet := range xlsxFile.Sheets {
    fmt.Printf("%d: %s\n", i + 1, sheet.Name)
  }
}

func RunXLSX(args []string) {
  fs := flag.NewFlagSet("xlsx", flag.ExitOnError)
  var dirname, sheet string
  var listSheets bool
  fs.BoolVar(&listSheets, "list-sheets", false, "List sheets in file")
  fs.StringVar(&dirname, "dirname", "", "Name of folder to output sheets to")
  fs.StringVar(&sheet, "sheet", "", "Name of sheet to convert")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only convert one file")
    os.Exit(1)
  } else if len(moreArgs) < 1 {
    fmt.Fprintln(os.Stderr, "Cannot convert file from stdin")
    os.Exit(1)
  }
  filename := moreArgs[0]
  if listSheets {
    ListXLXSSheets(filename)
  } else {
    if sheet == "" {
      if dirname == "" {
        fileParts := strings.Split(filename, ".")
        dirname = strings.Join(fileParts[:len(fileParts) - 1], ".")
      }
      ConvertXLSXFull(filename, dirname)
    } else {
      ConvertXLSXSheet(filename, sheet)
    }
  }
}
