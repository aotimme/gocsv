package main

import (
  "encoding/csv"
  "fmt"
  "os"
)


func ShowHeaders(reader *csv.Reader) {
  header, err := reader.Read()
  if err != nil {
    panic(err)
  }
  for i, name := range header {
    fmt.Printf("%d: %s\n", i + 1, name)
  }
}


func RunHeaders(args []string) {
  if len(args) > 1 {
    fmt.Fprintln(os.Stderr, "Can only show headers for one table")
    os.Exit(1)
  }
  var reader *csv.Reader
  if len(args) == 1 {
    file, err := os.Open(args[0])
    if err != nil {
      panic(err)
    }
    defer file.Close()
    reader = csv.NewReader(file)
  } else {
    reader = csv.NewReader(os.Stdin)
  }
  ShowHeaders(reader)
}
