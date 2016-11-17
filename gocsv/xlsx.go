package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "os"
  "strings"
  "github.com/tealeg/xlsx"
)

func ConvertXLSXSheet(dirname string, sheet *xlsx.Sheet) {
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

func ConvertXLSX(filename, dirname string) {
  xlsxFile, err := xlsx.OpenFile(filename)
  if err != nil {
    panic(err)
  }
  err = os.Mkdir(dirname, os.ModeDir | 0755)
  for _, sheet := range xlsxFile.Sheets {
    ConvertXLSXSheet(dirname, sheet)
  }
}


func RunXLSX(args []string) {
  fs := flag.NewFlagSet("convert", flag.ExitOnError)
  var dirname string
  fs.StringVar(&dirname, "dir", "", "Name of folder sheets to")
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
  if dirname == "" {
    fileParts := strings.Split(filename, ".")
    dirname = strings.Join(fileParts[:len(fileParts) - 1], ".")
  }
  ConvertXLSX(filename, dirname)
}
