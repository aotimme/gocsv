package cmd

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
	useTimeLayoutEnvVar()

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	Stats(inputCsvs[0])
}

func Stats(inputCsv *InputCsv) {
	imc := NewInMemoryCsvFromInputCsv(inputCsv)
	imc.PrintStats()
}
