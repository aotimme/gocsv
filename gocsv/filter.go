package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
  "regexp"
)

func FilterRegex(inreader io.Reader, columns []string, expr string, exclude bool) {
  re, err := regexp.Compile(expr)
  if err != nil {
    panic(err)
  }
  reader := csv.NewReader(inreader)
  writer := csv.NewWriter(os.Stdout)

  // Read header to get column index and write.
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }
  columnIndices := make([]int, len(columns))
  for i, column := range columns {
    index := GetColumnIndexOrPanic(header, column)
    columnIndices[i] = index
  }

  writer.Write(header)
  writer.Flush()

  // Write filtered rows.
  for {
    row, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        break
      } else {
        panic(err)
      }
    }
    rowMatches := false
    for _, columnIndex := range columnIndices {
      if re.MatchString(row[columnIndex]) {
        rowMatches = true
        break
      }
    }
    shouldOutputRow := (exclude && !rowMatches) || (!exclude && rowMatches)
    if shouldOutputRow {
      writer.Write(row)
      writer.Flush()
    }
  }
}


func RunFilter(args []string) {
  fs := flag.NewFlagSet("filter", flag.ExitOnError)
  var regex, columnsString string
  var exclude bool
  fs.StringVar(&regex, "regex", "", "Regular expression for filtering")
  fs.StringVar(&columnsString, "columns", "", "Columns to filter against")
  fs.BoolVar(&exclude, "exclude", false, "Exclude matching rows")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  if regex == "" || columnsString == "" {
    fmt.Fprintln(os.Stderr, "Missing required arguments --regex or --columns\n")
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
  FilterRegex(inreader, columns, regex, exclude)
}
