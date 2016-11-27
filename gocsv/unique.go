package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
  "strconv"

  "github.com/alphagov/router/trie"
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


func UniqueifySortedWithCount(reader *csv.Reader, columns []string) {
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }

  shellRow := make([]string, len(header) + 1)

  columnIndices := GetColumnIndicesOrAll(columns, header)

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  copy(shellRow, header)
  shellRow[len(shellRow) - 1] = "Count"
  writer.Write(shellRow)
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
  numInRun := 1

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
    if rowMatchesOnIndices(row, lastRow, columnIndices) {
      numInRun++
    } else {
      copy(shellRow, lastRow)
      shellRow[len(shellRow) - 1] = strconv.Itoa(numInRun)
      writer.Write(shellRow)
      writer.Flush()
      lastRow = row
      numInRun = 1
    }
  }
  copy(shellRow, lastRow)
  shellRow[len(shellRow) - 1] = strconv.Itoa(numInRun)
  writer.Write(shellRow)
  writer.Flush()
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

  seenRowsTrie := trie.NewTrie()
  lastRowArray := make([]string, len(columnIndices))

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
    for i, columnIndex := range columnIndices {
      lastRowArray[i] = row[columnIndex]
    }
    _, ok := seenRowsTrie.Get(lastRowArray)
    if !ok {
      seenRowsTrie.Set(lastRowArray, true)
      writer.Write(row)
      writer.Flush()
    }
  }
}


func UniqueifyUnsortedWithCount(reader *csv.Reader, columns []string) {
  imc := NewInMemoryCsv(reader)

  columnIndices := GetColumnIndicesOrAll(columns, imc.header)

  rowIndexToCount := make(map[int]int)
  seenRowsTrie := trie.NewTrie()

  lastRowArray := make([]string, len(columnIndices))

  for rowIndex, row := range imc.rows {
    for i, columnIndex := range columnIndices {
      lastRowArray[i] = row[columnIndex]
    }
    val, ok := seenRowsTrie.Get(lastRowArray)
    if ok {
      previousRowIndex := val.(int)
      rowIndexToCount[previousRowIndex] = rowIndexToCount[previousRowIndex] + 1
    } else {
      previousRowIndex := rowIndex
      seenRowsTrie.Set(lastRowArray, previousRowIndex)
      rowIndexToCount[previousRowIndex] = 1
    }
  }

  shellRow := make([]string, len(imc.header) + 1)
  copy(shellRow, imc.header)
  shellRow[len(shellRow) - 1] = "Count"

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  writer.Write(shellRow)
  writer.Flush()

  // Write unique rows with count.
  for rowIndex, row := range imc.rows {
    count, ok := rowIndexToCount[rowIndex]
    if ok {
      copy(shellRow, row)
      shellRow[len(shellRow) - 1] = strconv.Itoa(count)
      writer.Write(shellRow)
      writer.Flush()
    }
  }
}

func RunUnique(args []string) {
  fs := flag.NewFlagSet("unique", flag.ExitOnError)
  var columnsString string
  var sorted, count bool
  fs.StringVar(&columnsString, "columns", "", "Columns to use for comparison")
  fs.StringVar(&columnsString, "c", "", "Columns to use for comparison (shorthand)")
  fs.BoolVar(&sorted, "sorted", false, "Whether input CSV is already sorted")
  fs.BoolVar(&count, "count", false, "Whether to append a Count column")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
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
    if count {
      UniqueifySortedWithCount(reader, columns)
    } else {
      UniqueifySorted(reader, columns)
    }
  } else {
    if count {
      UniqueifyUnsortedWithCount(reader, columns)
    } else {
      UniqueifyUnsorted(reader, columns)
    }
  }
}
