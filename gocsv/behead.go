package main

import (
  "encoding/csv"
  "fmt"
  "io"
  "os"
)

func Behead(reader *csv.Reader) {
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
  Behead(reader)
}
