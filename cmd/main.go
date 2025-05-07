package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"text/tabwriter"
	"time"
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
	usage += version() + "\n"
	usage += "Subcommands:\n"
	for _, subcommand := range subcommands {
		usage += usageForSubcommand(subcommand)
	}
	usage += "See https://github.com/aotimme/gocsv for more documentation."
	return usage
}

func version() string {
	if VERSION != "" && GIT_HASH != "" {
		return fmt.Sprintf("Version: %s (%s)", VERSION, GIT_HASH)
	}

	s := ""
	if bi, ok := debug.ReadBuildInfo(); ok {
		s += "go:\t" + bi.GoVersion + "\n"
		for _, x := range bi.Settings {
			if x.Key == "vcs.revision" {
				// short hash
				s += "vcs.revision:\t" + x.Value[:7] + "\n"
			}
			if x.Key == "vcs.time" {
				t, _ := time.Parse(time.RFC3339, x.Value)
				t = t.Local()
				s += "vcs.time:\t" + t.Format(time.RFC3339) + "\n"
			}
			if x.Key == "vcs.modified" {
				s += "vcs.modified:\t" + x.Value + "\n"
			}
		}
	}
	s += "local-build:\t" + time.Now().Format(time.RFC3339)

	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	fmt.Fprint(w, s)
	w.Flush()

	return buf.String()
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
		fmt.Println(version())
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
