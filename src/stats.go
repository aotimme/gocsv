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
func (sub *StatsSubcommand) SetFlags(fs *flag.FlagSet) {
}

func (sub *StatsSubcommand) Run(args []string) {
	inputCsvs, err := GetInputCsvs(args, 1)
	if err != nil {
		panic(err)
	}

	Stats(inputCsvs[0])
}

func Stats(inputCsv AbstractInputCsv) {
	imc := NewInMemoryCsvFromInputCsv(inputCsv)
	imc.PrintStats()
}
