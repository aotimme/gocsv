package main

import (
  "encoding/csv"
  "fmt"
  "io"
  "os"
)

func Tsv(inreader io.Reader) {
  reader := csv.NewReader(inreader)
  writer := csv.NewWriter(os.Stdout)
  writer.Comma = '\t'

  // Write all rows with tabs.
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


func RunTsv(args []string) {
  if len(args) > 1 {
    fmt.Fprintln(os.Stderr, "Can only convert one table to TSV")
    os.Exit(2)
  }
  var inreader io.Reader
  if len(args) == 1 {
    file, err := os.Open(args[0])
    if err != nil {
      panic(err)
    }
    defer file.Close()
    inreader = file
  } else {
    inreader = os.Stdin
  }
  Tsv(inreader)
}
