package main

import (
	"fmt"
	"os"
)

// Keep this in sync with the README.
func usage() string {
	return `Usage:
  Valid subcommands are:
  - describe
      Get basic information about a CSV.
  - dimensions (alias: dims)
      Get the dimensions of a CSV.
  - clean
      Clean a CSV of common formatting issues.
  - delimiter (alias: delim)
      Change the delimiter being used for a CSV.
  - tsv
      Transform a CSV into a TSV.
  - headers
      View the headers from a CSV.
  - view
      Display a CSV in a pretty tabular format.
  - stats
      Get some basic statistics on a CSV.
  - rename
      Rename the headers of a CSV.
  - head
      Extract the first N rows from a CSV.
  - tail
      Extract the last N rows from a CSV.
  - behead
      Remove header row(s) from a CSV.
  - autoincrement (alias: autoinc)
      Add a column of incrementing integers to a CSV.
  - stack
      Stack multiple CSVs into one CSV.
  - split
      Split a CSV into multiple files.
  - sort
      Sort a CSV based on one or more columns.
  - filter
      Extract rows whose column match some criterion.
  - replace
      Replace values in cells by regular expression.
  - select
      Extract specified columns.
  - sample
      Sample rows.
  - unique (alias: uniq)
      Extract unique rows based upon certain columns.
  - join
      Join two CSVs based on equality of elements in a column.
  - xlsx
      Convert sheets of a XLSX file to CSV.
See https://github.com/DataFoxCo/gocsv for more documentation.`
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Fprintln(os.Stderr, "Must provide a valid subcommand.")
		fmt.Fprintf(os.Stderr, "%s\n", usage())
		os.Exit(1)
		return
	}
	subcommand := args[1]
	if subcommand == "describe" {
		RunDescribe(args[2:])
	} else if subcommand == "dimensions" || subcommand == "dims" {
		RunDimensions(args[2:])
	} else if subcommand == "headers" {
		RunHeaders(args[2:])
	} else if subcommand == "view" {
		RunView(args[2:])
	} else if subcommand == "stats" {
		RunStats(args[2:])
	} else if subcommand == "rename" {
		RunRename(args[2:])
	} else if subcommand == "clean" {
		RunClean(args[2:])
	} else if subcommand == "tsv" {
		RunTsv(args[2:])
	} else if subcommand == "delim" || subcommand == "delimiter" {
		RunDelimiter(args[2:])
	} else if subcommand == "head" {
		RunHead(args[2:])
	} else if subcommand == "tail" {
		RunTail(args[2:])
	} else if subcommand == "behead" {
		RunBehead(args[2:])
	} else if subcommand == "autoinc" || subcommand == "autoincrement" {
		RunAutoIncrement(args[2:])
	} else if subcommand == "stack" {
		RunStack(args[2:])
	} else if subcommand == "split" {
		RunSplit(args[2:])
	} else if subcommand == "filter" {
		RunFilter(args[2:])
	} else if subcommand == "replace" {
		RunReplace(args[2:])
	} else if subcommand == "select" {
		RunSelect(args[2:])
	} else if subcommand == "sort" {
		RunSort(args[2:])
	} else if subcommand == "sample" {
		RunSample(args[2:])
	} else if subcommand == "unique" || subcommand == "uniq" {
		RunUnique(args[2:])
	} else if subcommand == "join" {
		RunJoin(args[2:])
	} else if subcommand == "xlsx" {
		RunXLSX(args[2:])
	} else if subcommand == "help" {
		fmt.Fprintf(os.Stderr, "%s\n", usage())
	} else {
		fmt.Fprintf(os.Stderr, "Invalid subcommand \"%s\"\n", subcommand)
		fmt.Fprintf(os.Stderr, "%s\n", usage())
		os.Exit(1)
	}
}
