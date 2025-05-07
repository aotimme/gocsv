package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	// GIT_HASH is set during the build process using the -ldflags option.
	GIT_HASH string
	// VERSION is set during the build process using the -ldflags option.
	VERSION string
	// DEBUG is set by the common --debug flag
	DEBUG bool
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
	RegisterSubcommand(&AddSubcommand{})
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
	RegisterSubcommand(&NcolSubcommand{})
	RegisterSubcommand(&NrowSubcommand{})
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
	RegisterSubcommand(&TransposeSubcommand{})
	RegisterSubcommand(&TsvSubcommand{})
	RegisterSubcommand(&UniqueSubcommand{})
	RegisterSubcommand(&ViewSubcommand{})
	RegisterSubcommand(&XlsxSubcommand{})
	RegisterSubcommand(&ZipSubcommand{})
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

func exitTerseUsage() {
	fmt.Fprintln(os.Stderr, "usage: gocsv subcommand [flags] [file(s)]")
	fmt.Fprintln(os.Stderr, "Run 'gocsv help' for a list of subcommands, and 'gocsv <subcommand> -h' for subcommand details.")
	os.Exit(2)
}

// Keep this in sync with the README.
func usage() string {
	usage := "GoCSV is a command line CSV processing tool.\n"
	usage += fmt.Sprintf("Version: %s (%s)\n", VERSION, GIT_HASH)
	usage += "Subcommands:\n"
	for _, subcommand := range subcommands {
		usage += usageForSubcommand(subcommand)
	}
	usage += "See https://github.com/aotimme/gocsv for more documentation."
	return usage
}

func Main() {
	args := os.Args
	if len(args) == 1 {
		exitTerseUsage()
	}

	subcommandName := args[1]

	switch subcommandName {
	case "help":
		fmt.Println(usage())
		return
	case "version":
		fmt.Printf("%s (%s)\n", VERSION, GIT_HASH)
		return
	}

	for _, subcommand := range subcommands {
		if MatchesSubcommand(subcommand, subcommandName) {
			fs := flag.NewFlagSet(subcommand.Name(), flag.ExitOnError)
			fs.BoolVar(&DEBUG, "debug", false, "Enable debug mode")
			subcommand.SetFlags(fs)
			err := fs.Parse(args[2:])
			if err != nil {
				ExitWithError(err)
			}
			subcommand.Run(fs.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "error: invalid subcommand \"%s\"\n\n", subcommandName)
	exitTerseUsage()
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
