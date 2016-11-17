package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "os"
)


func SortCsv(reader *csv.Reader, columns []string, reverse, noInference bool) {
  imc := NewInMemoryCsv(reader)
  columnIndices := make([]int, len(columns))
  for i, column := range columns {
    columnIndices[i] = GetColumnIndexOrPanic(imc.header, column)
  }
  columnTypes := make([]ColumnType, len(columnIndices))
  for i, columnIndex := range columnIndices {
    if noInference {
      columnTypes[i] = STRING_TYPE
    } else {
      columnTypes[i] = imc.InferType(columnIndex)
    }
  }
  imc.SortRows(columnIndices, columnTypes, reverse)

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  writer.Write(imc.header)
  writer.Flush()

  // Write sorted rows.
  for _, row := range imc.rows {
    writer.Write(row)
    writer.Flush()
  }
}


func RunSort(args []string) {
  fs := flag.NewFlagSet("sort", flag.ExitOnError)
  var columnsString string
  var reverse, noInference bool
  fs.StringVar(&columnsString, "columns", "", "Columns to sort on")
  fs.StringVar(&columnsString, "c", "", "Columns to sort on (shorthand)")
  fs.BoolVar(&reverse, "reverse", false, "Sort in reverse")
  fs.BoolVar(&noInference, "no-inference", false, "Skip inference of input")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  if columnsString == "" {
    fmt.Fprintln(os.Stderr, "Missing argument --columns")
    os.Exit(1)
  }
  columns := GetArrayFromCsvString(columnsString)

  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only filter one table")
    os.Exit(1)
  }
  var reader *csv.Reader
  if len(moreArgs) == 1 {
    file, err := os.Open(moreArgs[0])
    if err != nil {
      panic(err)
    }
    defer file.Close()
    reader = csv.NewReader(file)
  } else {
    reader = csv.NewReader(os.Stdin)
  }

  SortCsv(reader, columns, reverse, noInference)
}
