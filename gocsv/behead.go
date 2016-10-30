package main

import (
  "encoding/csv"
  "fmt"
  "io"
  "os"
)

func Behead(inreader io.Reader) {
  reader := csv.NewReader(inreader)
  writer := csv.NewWriter(os.Stdout)

  // Get rid of the header.
  _, err := reader.Read()
  if err != nil {
    panic(err)
  }

  // Write rows.
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


func RunBehead(args []string) {
  if len(args) > 1 {
    fmt.Fprintln(os.Stderr, "Can only behead one table")
    return
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
  Behead(inreader)
}
