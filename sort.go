package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
)


func SortCsv(inreader io.Reader, columns []string, reverse, noInference bool) {
  imc := NewInMemoryCsv(inreader)
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
  fs := flag.NewFlagSet("sort", flag.PanicOnError)
  var columnsString string
  var reverse, noInference bool
  fs.StringVar(&columnsString, "columns", "", "Columns to sort on")
  fs.BoolVar(&reverse, "reverse", false, "Sort in reverse")
  fs.BoolVar(&noInference, "no-inference", false, "Skip inference of input")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  if columnsString == "" {
    fmt.Fprintln(os.Stderr, "Missing argument --columns")
    return
  }
  columns := GetArrayFromCsvString(columnsString)

  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only filter one table")
    return
  }
  var inreader io.Reader
  if len(moreArgs) == 1 {
    file, err := os.Open(moreArgs[0])
    if err != nil {
      panic(err)
    }
    defer file.Close()
    inreader = file
  } else {
    inreader = os.Stdin
  }

  SortCsv(inreader, columns, reverse, noInference)
}
