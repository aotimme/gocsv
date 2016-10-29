package main

import (
  "fmt"
  "os"
)


func main() {
  args := os.Args
  if len(args) == 1 {
    fmt.Fprintln(os.Stderr, "Valid subcommands are \"filter\" or \"select\"")
    return
  }
  subcommand := args[1]
  if subcommand == "headers" {
    RunHeaders(args[2:])
  } else if subcommand == "behead" {
    RunBehead(args[2:])
  } else if subcommand == "stack" {
    RunStack(args[2:])
  } else if subcommand == "filter" {
    RunFilter(args[2:])
  } else if subcommand == "select" {
    RunSelect(args[2:])
  } else if subcommand == "sort" {
    RunSort(args[2:])
  } else if subcommand == "join" {
    RunJoin(args[2:])
  } else {
    fmt.Fprintf(os.Stderr, "Invalid subcommand \"\"\n", subcommand)
    return
  }
}
