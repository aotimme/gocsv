package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
  "strconv"
)

func AutoIncrement(reader *csv.Reader, name string, seed int, prepend bool) {
  writer := csv.NewWriter(os.Stdout)

  // Read and write header.
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }
  numInputColumns := len(header)
  shellRow := make([]string, numInputColumns + 1)
  if prepend {
    shellRow[0] = name
    for i, elem := range header {
      shellRow[i + 1] = elem
    }
  } else {
    copy(shellRow, header)
    shellRow[numInputColumns] = name
  }
  writer.Write(shellRow)
  writer.Flush()

  // Write rows with autoincrement.
  inc := seed
  for {
    row, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        break
      } else {
        panic(err)
      }
    }
    incStr := strconv.Itoa(inc)
    if prepend {
      shellRow[0] = incStr
      for i, elem := range row {
        shellRow[i + 1] = elem
      }
    } else {
      copy(shellRow, row)
      shellRow[numInputColumns] = incStr
    }
    inc++
    writer.Write(shellRow)
    writer.Flush()
  }
}


func RunAutoIncrement(args []string) {
  fs := flag.NewFlagSet("autoincrement", flag.ExitOnError)
  var name string
  var seed int
  var prepend bool
  fs.StringVar(&name, "name", "ID", "Name of autoincrementing column")
  fs.IntVar(&seed, "seed", 1, "Initial value of autoincrementing column")
  fs.BoolVar(&prepend, "prepend", false, "Prepend the autoincrementing column (defaults to append)")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only autoincrement one file")
    return
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
  AutoIncrement(reader, name, seed, prepend)
}
