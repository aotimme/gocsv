package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
  "regexp"
)

func FilterRegex(inreader io.Reader, column, expr string) {
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
  columnIndex := GetColumnIndexOrPanic(header, column)

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
    if re.MatchString(row[columnIndex]) {
      writer.Write(row)
      writer.Flush()
    }
  }
}


func RunFilter(args []string) {
  fs := flag.NewFlagSet("filter", flag.PanicOnError)
  var regex, column string
  fs.StringVar(&regex, "regex", "", "Regular expression for filtering")
  fs.StringVar(&column, "column", "", "Column to filter against")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  if regex == "" || column == "" {
    fmt.Fprintln(os.Stderr, "Missing required arguments --regex or --column\n")
    return
  }
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
  FilterRegex(inreader, column, regex)
}
