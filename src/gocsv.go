package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Subcommand interface {
	Name() string
	Aliases() []string
	Description() string
	SetFlags(*flag.FlagSet)
	Run([]string)
}

var subcommands []Subcommand

func RegisterSubcommand(sub Subcommand) {
	subcommands = append(subcommands, sub)
}

func init() {
	RegisterSubcommand(&AutoincrementSubcommand{})
	RegisterSubcommand(&BeheadSubcommand{})
	RegisterSubcommand(&CapSubcommand{})
	RegisterSubcommand(&CleanSubcommand{})
	RegisterSubcommand(&DelimiterSubcommand{})
	RegisterSubcommand(&DescribeSubcommand{})
	RegisterSubcommand(&DimensionsSubcommand{})
	RegisterSubcommand(&FilterSubcommand{})
	RegisterSubcommand(&HeadSubcommand{})
	RegisterSubcommand(&HeadersSubcommand{})
	RegisterSubcommand(&JoinSubcommand{})
	RegisterSubcommand(&RenameSubcommand{})
	RegisterSubcommand(&ReplaceSubcommand{})
	RegisterSubcommand(&SampleSubcommand{})
	RegisterSubcommand(&SelectSubcommand{})
	RegisterSubcommand(&SortSubcommand{})
	RegisterSubcommand(&SplitSubcommand{})
	RegisterSubcommand(&SqlSubcommand{})
	RegisterSubcommand(&StackSubcommand{})
	RegisterSubcommand(&StatsSubcommand{})
	RegisterSubcommand(&TailSubcommand{})
	RegisterSubcommand(&TsvSubcommand{})
	RegisterSubcommand(&UniqueSubcommand{})
	RegisterSubcommand(&ViewSubcommand{})
	RegisterSubcommand(&XlsxSubcommand{})
}

func usageForSubcommand(subcommand Subcommand) string {
	retval := "  - " + subcommand.Name()
	aliases := subcommand.Aliases()
	if len(aliases) == 1 {
		retval += fmt.Sprintf(" (alias: %s)", aliases[0])
	} else if len(aliases) > 1 {
		retval += fmt.Sprintf(" (aliases: %s)", strings.Join(aliases, ", "))
	}
	retval += fmt.Sprintf("\n      %s\n", subcommand.Description())
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
		if MatchesSubcommand(subcommand, subcommandName) {
			fs := flag.NewFlagSet(subcommand.Name(), flag.ExitOnError)
			subcommand.SetFlags(fs)
			err := fs.Parse(args[2:])
			if err != nil {
				panic(err)
			}
			subcommand.Run(fs.Args())
			return
		}
	}
	fmt.Fprintf(os.Stderr, "Invalid subcommand \"%s\"\n", subcommandName)
	fmt.Fprintf(os.Stderr, "%s\n", usage())
	os.Exit(1)
}

func MatchesSubcommand(sub Subcommand, name string) bool {
	if name == sub.Name() {
		return true
	}
	for _, alias := range sub.Aliases() {
		if alias == name {
			return true
		}
	}
	return false
}
