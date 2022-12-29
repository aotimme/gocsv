package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type SelectSubcommand struct {
	columnsString string
	exclude       bool
	rawOutput     bool
}

func (sub *SelectSubcommand) Name() string {
	return "select"
}
func (sub *SelectSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SelectSubcommand) Description() string {
	return "Extract specified columns."
}
func (sub *SelectSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.columnsString, "columns", "", "Columns to select")
	fs.StringVar(&sub.columnsString, "c", "", "Columns to select (shorthand)")
	fs.BoolVar(&sub.exclude, "exclude", false, "Whether to exclude the specified columns")
	fs.BoolVar(&sub.rawOutput, "raw-output", false, "Whether to output as raw lines (no CSV formatting) -- only applies if select returns one column")
	fs.BoolVar(&sub.rawOutput, "r", false, "Whether to output as raw lines (no CSV formatting) -- only applies if select returns one column (shorthand)")
}

func (sub *SelectSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsvs(inputCsvs)
	outputCsv.SetWriteRaw(sub.rawOutput)
	sub.RunSelect(inputCsvs[0], outputCsv)
}

func (sub *SelectSubcommand) RunSelect(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	if sub.columnsString == "" {
		fmt.Fprintln(os.Stderr, "Missing required argument --columns")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(sub.columnsString)

	if sub.exclude {
		ExcludeColumns(inputCsv, outputCsvWriter, columns)
	} else {
		SelectColumns(inputCsv, outputCsvWriter, columns)
	}
}

func ExcludeColumns(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, columns []string) {
	// Get the column indices to exclude.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	columnIndices := GetIndicesForColumnsOrPanic(header, columns)
	columnIndicesToExclude := make(map[int]bool)
	for _, columnIndex := range columnIndices {
		columnIndicesToExclude[columnIndex] = true
	}

	outrow := make([]string, len(header)-len(columnIndicesToExclude))

	// Write header
	curIdx := 0
	for index, elem := range header {
		_, exclude := columnIndicesToExclude[index]
		if !exclude {
			outrow[curIdx] = elem
			curIdx++
		}
	}

	outputCsvWriter.Write(outrow)

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		curIdx = 0
		for index, elem := range row {
			_, exclude := columnIndicesToExclude[index]
			if !exclude {
				outrow[curIdx] = elem
				curIdx++
			}
		}
		outputCsvWriter.Write(outrow)
	}
}

func SelectColumns(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, columns []string) {
	// Get the column indices to write.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	columnIndices := GetIndicesForColumnsOrPanic(header, columns)
	outrow := make([]string, len(columnIndices))
	for i, columnIndex := range columnIndices {
		outrow[i] = header[columnIndex]
	}
	outputCsvWriter.Write(outrow)

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		for i, columnIndex := range columnIndices {
			outrow[i] = row[columnIndex]
		}
		outputCsvWriter.Write(outrow)
	}
}
