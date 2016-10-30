package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
  "strings"
)


func concat(outrow, row1, row2 []string) {
  i := 0
  for _, elem := range row1 {
    outrow[i] = elem
    i++
  }
  for _, elem := range row2 {
    outrow[i] = elem
    i++
  }
}


func InnerJoin(leftReader, rightReader io.Reader, leftColname, rightColname string) {
  leftCsv := NewInMemoryCsv(leftReader)
  leftColIndex := GetColumnIndexOrPanic(leftCsv.header, leftColname)
  rightCsv := NewInMemoryCsv(rightReader)
  rightColIndex := GetColumnIndexOrPanic(rightCsv.header, rightColname)
  rightCsv.Index(rightColIndex)


  numLeftColumns := len(leftCsv.header)
  numRightColumns := len(rightCsv.header)

  shellRow := make([]string, numLeftColumns + numRightColumns)

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  concat(shellRow, leftCsv.header, rightCsv.header)
  writer.Write(shellRow)
  writer.Flush()

  // Write inner-joined rows.
  for _, row := range leftCsv.rows {
    rightRows := rightCsv.GetRowsMatchingIndexedColumn(row[leftColIndex])
    if len(rightRows) > 0 {
      for _, rightRow := range rightRows {
        concat(shellRow, row, rightRow)
        writer.Write(shellRow)
        writer.Flush()
      }
    }
  }
}


func LeftJoin(leftReader, rightReader io.Reader, leftColname, rightColname string) {
  leftCsv := NewInMemoryCsv(leftReader)
  leftColIndex := GetColumnIndexOrPanic(leftCsv.header, leftColname)

  rightCsv := NewInMemoryCsv(rightReader)
  rightColIndex := GetColumnIndexOrPanic(rightCsv.header, rightColname)
  rightCsv.Index(rightColIndex)

  numLeftColumns := len(leftCsv.header)
  numRightColumns := len(rightCsv.header)

  emptyRightRow := make([]string, numRightColumns)
  shellRow := make([]string, numLeftColumns + numRightColumns)

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  concat(shellRow, leftCsv.header, rightCsv.header)
  writer.Write(shellRow)
  writer.Flush()

  // Write right-joined rows.
  for _, row := range leftCsv.rows {
    rightRows := rightCsv.GetRowsMatchingIndexedColumn(row[leftColIndex])
    if len(rightRows) > 0 {
      for _, rightRow := range rightRows {
        concat(shellRow, row, rightRow)
        writer.Write(shellRow)
        writer.Flush()
      }
    } else {
      concat(shellRow, row, emptyRightRow)
      writer.Write(shellRow)
      writer.Flush()
    }
  }
}


func RightJoin(leftReader, rightReader io.Reader, leftColname, rightColname string) {
  leftCsv := NewInMemoryCsv(leftReader)
  leftColIndex := GetColumnIndexOrPanic(leftCsv.header, leftColname)
  leftCsv.Index(leftColIndex)

  rightCsv := NewInMemoryCsv(rightReader)
  rightColIndex := GetColumnIndexOrPanic(rightCsv.header, rightColname)

  numLeftColumns := len(leftCsv.header)
  numRightColumns := len(rightCsv.header)

  emptyLeftRow := make([]string, numLeftColumns)
  shellRow := make([]string, numLeftColumns + numRightColumns)

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  concat(shellRow, leftCsv.header, rightCsv.header)
  writer.Write(shellRow)
  writer.Flush()

  // Write left-joined rows.
  for _, row := range rightCsv.rows {
    leftRows := leftCsv.GetRowsMatchingIndexedColumn(row[rightColIndex])
    if len(leftRows) > 0 {
      for _, leftRow := range leftRows {
        concat(shellRow, leftRow, row)
        writer.Write(shellRow)
        writer.Flush()
      }
    } else {
      concat(shellRow, emptyLeftRow, row)
      writer.Write(shellRow)
      writer.Flush()
    }
  }
}


func OuterJoin(leftReader, rightReader io.Reader, leftColname, rightColname string) {
  // Basically do a left join and then append any rows from the right table
  // that weren't already included.

  leftCsv := NewInMemoryCsv(leftReader)
  leftColIndex := GetColumnIndexOrPanic(leftCsv.header, leftColname)

  rightCsv := NewInMemoryCsv(rightReader)
  rightColIndex := GetColumnIndexOrPanic(rightCsv.header, rightColname)
  rightCsv.Index(rightColIndex)


  numLeftColumns := len(leftCsv.header)
  numRightColumns := len(rightCsv.header)

  emptyLeftRow := make([]string, numLeftColumns)
  emptyRightRow := make([]string, numRightColumns)
  shellRow := make([]string, len(leftCsv.header) + len(rightCsv.header))

  // whether the row in the right column has been included already.
  rightIncludeStatus := make([]bool, len(rightCsv.rows))

  writer := csv.NewWriter(os.Stdout)

  // Write header.
  concat(shellRow, leftCsv.header, rightCsv.header)
  writer.Write(shellRow)
  writer.Flush()

  // Write left-joined rows.
  for _, row := range leftCsv.rows {
    rightRowIndices := rightCsv.GetRowIndicesMatchingIndexedColumn(row[leftColIndex])
    if len(rightRowIndices) > 0 {
      for _, rightRowIndex := range rightRowIndices {
        rightIncludeStatus[rightRowIndex] = true
        concat(shellRow, row, rightCsv.rows[rightRowIndex])
        writer.Write(shellRow)
        writer.Flush()
      }
    } else {
      concat(shellRow, row, emptyRightRow)
      writer.Write(shellRow)
      writer.Flush()
    }
  }

  // Write remaining right rows.
  for i, row := range rightCsv.rows {
    if rightIncludeStatus[i] {
      continue
    }
    concat(shellRow, emptyLeftRow, row)
    writer.Write(shellRow)
    writer.Flush()
  }
}


func RunJoin(args []string) {
  fs := flag.NewFlagSet("join", flag.ExitOnError)
  var columnsString string
  var left, right, outer bool
  fs.StringVar(&columnsString, "columns", "", "Columns to join on")
  fs.BoolVar(&left, "left", false, "Left join")
  fs.BoolVar(&right, "right", false, "Right join")
  fs.BoolVar(&outer, "outer", false, "Full outer join")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  numJoins := 0
  if left {
    numJoins++
  }
  if right {
    numJoins++
  }
  if outer {
    numJoins++
  }
  if numJoins > 1 {
    fmt.Fprintln(os.Stderr, "Must only specify zero or one of --left, --right, or --outer")
    return
  }
  c := csv.NewReader(strings.NewReader(columnsString))
  rows, err := c.ReadAll()
  if err != nil {
    panic(err)
  }
  columns := rows[0]
  if len(columns) < 1 || len(columns) > 2 {
    fmt.Fprintln(os.Stderr, "Invalid argument for --columns")
    return
  }
  if len(columns) == 1 {
    columns = append(columns, columns[0])
  }
  moreArgs := fs.Args()
  if len(moreArgs) == 0 {
    fmt.Fprintln(os.Stderr, "Missing right table to join against")
    return
  } else if len(moreArgs) > 2 {
    fmt.Fprintln(os.Stderr, "Too many join tables")
    return
  }
  var leftReader, rightReader io.Reader
  if len(moreArgs) == 1 {
    leftReader = os.Stdin
    file, err := os.Open(moreArgs[0])
    if err != nil {
      panic(err)
    }
    defer file.Close()
    rightReader = file
  } else {
    file, err := os.Open(moreArgs[0])
    if err != nil {
      panic(err)
    }
    defer file.Close()
    leftReader = file
    file, err = os.Open(moreArgs[1])
    if err != nil {
      panic(err)
    }
    defer file.Close()
    rightReader = file
  }
  if left {
    LeftJoin(leftReader, rightReader, columns[0], columns[1])
  } else if right {
    RightJoin(leftReader, rightReader, columns[0], columns[1])
  } else if outer {
    OuterJoin(leftReader, rightReader, columns[0], columns[1])
  } else {
    InnerJoin(leftReader, rightReader, columns[0], columns[1])
  }
}
