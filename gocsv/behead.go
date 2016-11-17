package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
)

func Behead(reader *csv.Reader, numHeaders int) {
  writer := csv.NewWriter(os.Stdout)

  // Get rid of the header rows.
  for i := 0; i < numHeaders; i++ {
    _, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        // If we remove _all_ the headers, then end early.
        return
      } else {
        panic(err)
      }
    }
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
  fs := flag.NewFlagSet("behead", flag.ExitOnError)
  var numHeaders int
  fs.IntVar(&numHeaders, "n", 1, "Number of headers to remove")

  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }

  if numHeaders < 1 {
    fmt.Fprintln(os.Stderr, "Invalid argument -n")
    os.Exit(1)
  }

  // Get input CSV
  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only behead one table")
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
  Behead(reader, numHeaders)
}
