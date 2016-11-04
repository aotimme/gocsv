package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "math"
  "os"
  "regexp"
  "strconv"
)

func FilterMatchFunc(reader *csv.Reader, columns []string, exclude bool, matchFunc func(string) bool ) {
  writer := csv.NewWriter(os.Stdout)

  // Read header to get column index and write.
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }

  // Get indices to compare against.
  // If no columns are specified, then check against all.
  var columnIndices []int
  if len(columns) == 0 {
    columnIndices = make([]int, len(header))
    for i, _ := range header {
      columnIndices[i] = i
    }
  } else {
    columnIndices = make([]int, len(columns))
    for i, column := range columns {
      index := GetColumnIndexOrPanic(header, column)
      columnIndices[i] = index
    }
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
      if matchFunc(row[columnIndex]) {
        rowMatches = true
        break
      }
    }
    shouldOutputRow := (!exclude && rowMatches) || (exclude && !rowMatches)
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
  var gt, gte, lt, lte float64
  positiveInfinity := math.Inf(1)
  negativeInfinity := math.Inf(-1)
  fs.StringVar(&regex, "regex", "", "Regular expression for filtering")
  fs.StringVar(&columnsString, "columns", "", "Columns to filter against")
  fs.BoolVar(&exclude, "exclude", false, "Exclude matching rows")
  fs.Float64Var(&gt, "gt", negativeInfinity, "Greater than")
  fs.Float64Var(&gte, "gte", negativeInfinity, "Greater than or equal to")
  fs.Float64Var(&lt, "lt", positiveInfinity, "Less than")
  fs.Float64Var(&lte, "lte", positiveInfinity, "Less than or equal to")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }

  // Get columns to compare against
  var columns []string
  if columnsString == "" {
    columns = make([]string, 0)
  } else {
    columns = GetArrayFromCsvString(columnsString)
  }

  // Get match function
  var matchFunc func(string) bool
  if regex != "" {
    re, err := regexp.Compile(regex)
    if err != nil {
      panic(err)
    }
    matchFunc = func(elem string) bool {
      return re.MatchString(elem)
    }
  } else if gt != negativeInfinity {
    matchFunc = func(elem string) bool {
      elem64, err := strconv.ParseFloat(elem, 64)
      if err != nil {
        return false
      }
      return elem64 > gt
    }
  } else if gte != negativeInfinity {
    matchFunc = func(elem string) bool {
      elem64, err := strconv.ParseFloat(elem, 64)
      if err != nil {
        return false
      }
      return elem64 >= gte
    }
  } else if lt != positiveInfinity {
    matchFunc = func(elem string) bool {
      elem64, err := strconv.ParseFloat(elem, 64)
      if err != nil {
        return false
      }
      return elem64 < lt
    }
  } else if lte != positiveInfinity {
    matchFunc = func(elem string) bool {
      elem64, err := strconv.ParseFloat(elem, 64)
      if err != nil {
        return false
      }
      return elem64 <= lte
    }
  } else {
    fmt.Fprintln(os.Stderr, "Missing filter function")
    os.Exit(2)
  }

  // Get input CSV
  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only filter one table")
    os.Exit(2)
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

  FilterMatchFunc(reader, columns, exclude, matchFunc)
}
