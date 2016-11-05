package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
)


func RenameColumns(reader *csv.Reader, columns, names []string) {
  writer := csv.NewWriter(os.Stdout)

  // Get the column indices to write.
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }
  renamedHeader := make([]string, len(header))
  copy(renamedHeader, header)

  for i, column := range columns {
    index := GetColumnIndexOrPanic(header, column)
    renamedHeader[index] = names[i]
  }

  writer.Write(renamedHeader)
  writer.Flush()

  for {
    row, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        break
      } else {
        panic(err)
      }
    }
    writer.Write(row)
    writer.Flush()
  }
}


func RunRename(args []string) {
  fs := flag.NewFlagSet("rename", flag.ExitOnError)
  var columnsString, namesString string
  fs.StringVar(&columnsString, "columns", "", "Columns to select")
  fs.StringVar(&namesString, "names", "", "New names for columns")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  if columnsString == "" {
    fmt.Fprintf(os.Stderr, "Missing required argument --columns")
    os.Exit(1)
  }
  if namesString == "" {
    fmt.Fprintf(os.Stderr, "Missing required argument --names")
    os.Exit(1)
  }
  columns := GetArrayFromCsvString(columnsString)
  names := GetArrayFromCsvString(namesString)
  if len(columns) != len(names) {
    fmt.Fprintln(os.Stderr, "Length of --columns and --names argument must be the same")
    os.Exit(1)
  }
  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only call rename on one table")
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
  RenameColumns(reader, columns, names)
}
