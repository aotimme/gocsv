package cmd

import "flag"

type TransposeSubcommand struct {
}

func (sub *TransposeSubcommand) Name() string {
	return "transpose"
}

func (sub *TransposeSubcommand) Aliases() []string {
	return []string{}
}

func (sub *TransposeSubcommand) Description() string {
	return "Transpose a CSV"
}

func (sub *TransposeSubcommand) SetFlags(fs *flag.FlagSet) {
}

func (sub *TransposeSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsvs(inputCsvs)
	sub.RunTranspose(inputCsvs[0], outputCsv)
}

func (sub *TransposeSubcommand) RunTranspose(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	imc := NewInMemoryCsvFromInputCsv(inputCsv)

	numRows := imc.NumRows()
	numColumns := imc.NumColumns()

	outrow := make([]string, numRows+1)
	for j := 0; j < numColumns; j++ {
		outrow[0] = imc.header[j]
		for i := 0; i < numRows; i++ {
			outrow[i+1] = imc.rows[i][j]
		}
		outputCsvWriter.Write(outrow)
	}
}
