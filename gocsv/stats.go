package main

import (
	"flag"
)

type StatsSubcommand struct{}

func (sub *StatsSubcommand) Name() string {
	return "stats"
}
func (sub *StatsSubcommand) Aliases() []string {
	return []string{}
}
func (sub *StatsSubcommand) Description() string {
	return "Get some basic statistics on a CSV."
}

func (sub *StatsSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	Stats(inputCsvs[0])
}

func Stats(inputCsv AbstractInputCsv) {
	imc := NewInMemoryCsvFromInputCsv(inputCsv)
	imc.PrintStats()
}
