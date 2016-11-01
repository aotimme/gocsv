package main

import (
  "fmt"
  "os"
)


// Keep this in sync with the README.
func usage() string {
  return `Usage:
  Valid subcommands are:
  - clean
    Clean a CSV of common formatting issues.
  - tsv
    Transform a CSV into a TSV.
  - headers
    View the headers from a CSV.
  - head
    Extract the first N rows from a CSV.
  - tail
    Extract the last N rows from a CSV.
  - behead
    Remove the header from a CSV.
  - autoincrement
    Add a column of incrementing integers to a CSV.
  - stack
    Stack multiple CSVs into one CSV.
  - split
    Split a CSV into multiple files.
  - sort
    Sort a CSV based on one or more columns.
  - filter
    Extract rows whose column matches a regular expression.
  - select
    Extract specified columns.
  - join
    Join two CSVs based on equality of elements in a column.
See https://github.com/DataFoxCo/gocsv for more documentation.`
}


func main() {
  args := os.Args
  if len(args) == 1 {
    fmt.Fprintln(os.Stderr, "Must provide a valid subcommand.")
    fmt.Fprintf(os.Stderr, "%s\n", usage())
    os.Exit(2)
    return
  }
  subcommand := args[1]
  if subcommand == "headers" {
    RunHeaders(args[2:])
  } else if subcommand == "clean" {
    RunClean(args[2:])
  } else if subcommand == "tsv" {
    RunTsv(args[2:])
  } else if subcommand == "head" {
    RunHead(args[2:])
  } else if subcommand == "tail" {
    RunTail(args[2:])
  } else if subcommand == "behead" {
    RunBehead(args[2:])
  } else if subcommand == "autoinc" || subcommand == "autoincrement"{
    RunAutoIncrement(args[2:])
  } else if subcommand == "stack" {
    RunStack(args[2:])
  } else if subcommand == "split" {
    RunSplit(args[2:])
  } else if subcommand == "filter" {
    RunFilter(args[2:])
  } else if subcommand == "select" {
    RunSelect(args[2:])
  } else if subcommand == "sort" {
    RunSort(args[2:])
  } else if subcommand == "join" {
    RunJoin(args[2:])
  } else if subcommand == "help" {
    fmt.Fprintf(os.Stderr, "%s\n", usage())
  } else {
    fmt.Fprintf(os.Stderr, "Invalid subcommand \"%s\"\n", subcommand)
    fmt.Fprintf(os.Stderr, "%s\n", usage())
    os.Exit(2)
  }
}
