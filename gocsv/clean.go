package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "os"
)

func Clean(reader *csv.Reader, noTrim bool) {
  writer := csv.NewWriter(os.Stdout)

  // Disable errors when fields are varying length
  reader.FieldsPerRecord = -1

  // Read in rows.
  rows, err := reader.ReadAll()
  if err != nil {
    panic(err)
  }

  // Determine how many columns there actually should be.
  numColumns := 0
  trimFromIndex := -1
  for i, row := range rows {
    lastNonEmptyIndex := -1
    for j, elem := range row {
      if elem != "" {
        lastNonEmptyIndex = j
      }
    }
    if lastNonEmptyIndex > -1 {
      trimFromIndex = -1
    } else if trimFromIndex == -1 {
      trimFromIndex = i
    }
    numColumnsInRow := lastNonEmptyIndex + 1
    if numColumns < numColumnsInRow {
      numColumns = numColumnsInRow
    }
  }


  // Fix rows and output them to writer.
  shellRow := make([]string, numColumns)
  for i, row := range rows {
    if !noTrim && trimFromIndex > -1 && i >= trimFromIndex {
      break
    }
    if len(row) == numColumns {
      // Just write the original row.
      writer.Write(row)
      writer.Flush()
    } else if len(row) < numColumns {
      // Pad the row.
      copy(shellRow, row)
      for i := len(row); i < numColumns; i++ {
        shellRow[i] = ""
      }
      writer.Write(shellRow)
      writer.Flush()
    } else {
      // Truncate the row.
      copy(shellRow, row)
      writer.Write(shellRow)
      writer.Flush()
    }
  }
}


func RunClean(args []string) {
  fs := flag.NewFlagSet("clean", flag.ExitOnError)
  var noTrim bool
  fs.BoolVar(&noTrim, "no-trim", false, "Don't trim end of file of empty rows")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  moreArgs := fs.Args()
  if len(moreArgs) > 1 {
    fmt.Fprintln(os.Stderr, "Can only clean one file")
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
  Clean(reader, noTrim)
}
