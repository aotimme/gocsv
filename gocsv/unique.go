package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
)


func GetColumnIndicesOrAll(columns, header []string) []int {
  var columnIndices []int
  if len(columns) == 0 {
    columnIndices = make([]int, len(header))
    for i := 0; i < len(header); i++ {
      columnIndices[i] = i
    }
  } else {
    columnIndices = make([]int, len(columns))
    for i, column := range columns {
      columnIndices[i] = GetColumnIndexOrPanic(header, column)
    }
  }
  return columnIndices
}

func rowMatchesOnIndices(rowA, rowB []string, columnIndices []int) bool {
  for _, columnIndex := range columnIndices {
    if rowA[columnIndex] != rowB[columnIndex] {
      return false
    }
  }
  return true
}

func UniqueifySorted(reader *csv.Reader, columns []string) {
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }

  columnIndices := GetColumnIndicesOrAll(columns, header)

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  writer.Write(header)
  writer.Flush()

  // Read and write first row.
  lastRow, err := reader.Read()
  if err != nil {
    if err == io.EOF {
      return
    } else {
      panic(err)
    }
  }
  writer.Write(lastRow)
  writer.Flush()

  // Write unique rows in order.
  for {
    row, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        break
      } else {
        panic(err)
      }
    }
    if !rowMatchesOnIndices(row, lastRow, columnIndices) {
      lastRow = row
      writer.Write(row)
      writer.Flush()
    }
  }
}

func isInIndexMap(row []string, indices []int, indexMap map[string]interface{}) bool {
  curMap := indexMap
  for _, index := range indices {
    cell := row[index]
    curMapInterface, ok := curMap[cell]
    if !ok {
      return false
    }
    curMap = curMapInterface.(map[string]interface{})
  }
  return true
}

func addToIndexMap(row []string, indices []int, indexMap map[string]interface{}) {
  prevMap := indexMap
  for _, index := range indices {
    cell := row[index]
    curMapInterface, ok := prevMap[cell]
    if ok {
      prevMap = curMapInterface.(map[string]interface{})
    } else {
      curMap := make(map[string]interface{})
      prevMap[cell] = curMap
      prevMap = curMap
    }
  }
}

func UniqueifyUnsorted(reader *csv.Reader, columns []string) {
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }

  columnIndices := GetColumnIndicesOrAll(columns, header)

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  writer.Write(header)
  writer.Flush()

  seenRowsIndexed := make(map[string]interface{})

  // Write unique rows in order.
  for {
    row, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        break
      } else {
        panic(err)
      }
    }

    if !isInIndexMap(row, columnIndices, seenRowsIndexed) {
      addToIndexMap(row, columnIndices, seenRowsIndexed)
      writer.Write(row)
      writer.Flush()
    }
  }
}

func RunUnique(args []string) {
  fs := flag.NewFlagSet("unique", flag.ExitOnError)
  var columnsString string
  var sorted bool
  fs.StringVar(&columnsString, "columns", "", "Columns to use for comparison")
  fs.StringVar(&columnsString, "c", "", "Columns to use for comparison (shorthand)")
  fs.BoolVar(&sorted, "sorted", false, "Whether input CSV is already sorted")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  if columnsString == "" {
    fmt.Fprintln(os.Stderr, "Missing argument --columns")
    os.Exit(1)
  }
  var columns []string
  if columnsString == "" {
    columns = make([]string, 0)
  } else {
    columns = GetArrayFromCsvString(columnsString)
  }

  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only unique one table")
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

  if sorted {
    UniqueifySorted(reader, columns)
  } else {
    UniqueifyUnsorted(reader, columns)
  }
}
