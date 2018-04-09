package main

import (
	"fmt"
	"os"
	"strings"
)

type Subcommand struct {
	name    string
	aliases []string
	usage   string
	runner  func([]string)
}

func (subcommand *Subcommand) Run(args []string) {
	subcommand.runner(args)
}

func (subcommand *Subcommand) DoesMatchName(name string) bool {
	if name == subcommand.name {
		return true
	}
	for _, alias := range subcommand.aliases {
		if alias == name {
			return true
		}
	}
	return false
}

var subcommands []Subcommand

func RegisterSubcommand(name string, aliases []string, usage string, runner func([]string)) {
	subcommand := Subcommand{
		name,
		aliases,
		usage,
		runner,
	}
	subcommands = append(subcommands, subcommand)
}

func init() {
	RegisterSubcommand("describe", []string{}, "Get basic information about a CSV.", RunDescribe)
	RegisterSubcommand("dimensions", []string{"dims"}, "Get the dimensions of a CSV.", RunDimensions)
	RegisterSubcommand("headers", []string{}, "View the headers from a CSV.", RunHeaders)
	RegisterSubcommand("view", []string{}, "Display a CSV in a pretty tabular format.", RunView)
	RegisterSubcommand("stats", []string{}, "Get some basic statistics on a CSV.", RunStats)
	RegisterSubcommand("rename", []string{}, "Rename the headers of a CSV.", RunRename)
	RegisterSubcommand("clean", []string{}, "Clean a CSV of common formatting issues.", RunClean)
	RegisterSubcommand("tsv", []string{}, "Transform a CSV into a TSV.", RunTsv)
	RegisterSubcommand("delimiter", []string{"delimiter"}, "Change the delimiter being used for a CSV.", RunDelimiter)
	RegisterSubcommand("head", []string{}, "Extract the first N rows from a CSV.", RunHead)
	RegisterSubcommand("tail", []string{}, "Extract the last N rows from a CSV.", RunTail)
	RegisterSubcommand("behead", []string{}, "Remove header row(s) from a CSV.", RunBehead)
	RegisterSubcommand("autoincrement", []string{"autoinc"}, "Add a column of incrementing integers to a CSV.", RunAutoIncrement)
	RegisterSubcommand("stack", []string{}, "Stack multiple CSVs into one CSV.", RunStack)
	RegisterSubcommand("split", []string{}, "Split a CSV into multiple files.", RunSplit)
	RegisterSubcommand("filter", []string{}, "Extract rows whose column match some criterion.", RunFilter)
	RegisterSubcommand("replace", []string{}, "Replace values in cells by regular expression.", RunReplace)
	RegisterSubcommand("select", []string{}, "Extract specified columns.", RunSelect)
	RegisterSubcommand("sort", []string{}, "Sort a CSV based on one or more columns.", RunSort)
	RegisterSubcommand("sample", []string{}, "Sample rows.", RunSample)
	RegisterSubcommand("unique", []string{"uniq"}, "Extract unique rows based upon certain columns.", RunUnique)
	RegisterSubcommand("join", []string{}, "Join two CSVs based on equality of elements in a column.", RunJoin)
	RegisterSubcommand("xlsx", []string{}, "Convert sheets of a XLSX file to CSV.", RunXLSX)
	RegisterSubcommand("sql", []string{}, "Run SQL queries on CSVs.", RunSql)
}

func usageForSubcommand(subcommand Subcommand) string {
	retval := "  - " + subcommand.name
	if len(subcommand.aliases) == 1 {
		retval += fmt.Sprintf(" (alias: %s)", subcommand.aliases[0])
	} else if len(subcommand.aliases) > 1 {
		retval += fmt.Sprintf(" (aliases: %s)", strings.Join(subcommand.aliases, ", "))
	}
	retval += fmt.Sprintf("\n      %s\n", subcommand.usage)
	return retval
}

// Keep this in sync with the README.
func usage() string {
	usage := "Usage:\n"
	usage += "  Valid subcommands are:\n"
	for _, subcommand := range subcommands {
		usage += usageForSubcommand(subcommand)
	}
	usage += "See https://github.com/DataFoxCo/gocsv for more documentation."
	return usage
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Fprintln(os.Stderr, "Must provide a valid subcommand.")
		fmt.Fprintf(os.Stderr, "%s\n", usage())
		os.Exit(1)
		return
	}
	subcommandName := args[1]
	if subcommandName == "help" {
		fmt.Fprintf(os.Stderr, "%s\n", usage())
		return
	}
	for _, subcommand := range subcommands {
		if subcommand.DoesMatchName(subcommandName) {
			subcommand.Run(args[2:])
			return
		}
	}
	fmt.Fprintf(os.Stderr, "Invalid subcommand \"%s\"\n", subcommandName)
	fmt.Fprintf(os.Stderr, "%s\n", usage())
	os.Exit(1)
}
