package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
  "regexp"
  "strconv"
  "strings"
)


func TailFromBottom(reader *csv.Reader, numRows int) {
  writer := csv.NewWriter(os.Stdout)

  // Read all rows.
  rows, err := reader.ReadAll()
  if err != nil {
    panic(err)
  }

  // Write header.
  writer.Write(rows[0])
  writer.Flush()

  // Write rows.
  startRow := len(rows) - numRows
  if startRow < 1 {
    startRow = 1
  }
  for i := startRow; i < len(rows); i++ {
    writer.Write(rows[i])
    writer.Flush()
  }
}


func TailFromTop(reader *csv.Reader, numRows int) {
  writer := csv.NewWriter(os.Stdout)

  // Read and write header.
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }
  writer.Write(header)
  writer.Flush()

  // Write rows after first `numRows` rows.
  curRow := 0
  for {
    row, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        break
      } else {
        panic(err)
      }
    }
    curRow++
    if curRow > numRows {
      writer.Write(row)
      writer.Flush()
    }
  }
}


func RunTail(args []string) {
  fs := flag.NewFlagSet("filter", flag.ExitOnError)
  var numRowsStr string
  fs.StringVar(&numRowsStr, "n", "10", "Number of rows to include")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  numRowsRegex := regexp.MustCompile("^\\+?\\d+$")
  if !numRowsRegex.MatchString(numRowsStr) {
    fmt.Fprintln(os.Stderr, "Invalid argument to -n")
    os.Exit(1)
  }
  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only run tail on one table")
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
  if strings.HasPrefix(numRowsStr, "+") {
    numRowsStr = strings.TrimPrefix(numRowsStr, "+")
    numRows, err := strconv.Atoi(numRowsStr)
    if err != nil {
      panic(err)
    }
    TailFromTop(reader, numRows)
  } else {
    numRows, err := strconv.Atoi(numRowsStr)
    if err != nil {
      panic(err)
    }
    TailFromBottom(reader, numRows)
  }
}
