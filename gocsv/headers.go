package main

import (
  "encoding/csv"
  "fmt"
  "io"
  "os"
)


func ShowHeaders(inreader io.Reader) {
  reader := csv.NewReader(inreader)
  for {
    header, err := reader.Read()
    if err != nil {
      panic(err)
    }
    for i, name := range header {
      fmt.Printf("%d: %s\n", i + 1, name)
    }
    break
  }
}


func RunHeaders(args []string) {
  if len(args) > 1 {
    fmt.Fprintln(os.Stderr, "Can only show headers for one table")
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
  ShowHeaders(inreader)
}
