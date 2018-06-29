package main

import (
	"flag"
	"io"
	"regexp"
)

type ReplaceSubcommand struct {
	columnsString   string
	regex           string
	repl            string
	caseInsensitive bool
}

func (sub *ReplaceSubcommand) Name() string {
	return "replace"
}
func (sub *ReplaceSubcommand) Aliases() []string {
	return []string{}
}
func (sub *ReplaceSubcommand) Description() string {
	return "Replace values in cells by regular expression."
}
func (sub *ReplaceSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.columnsString, "columns", "", "Columns to replace cells")
	fs.StringVar(&sub.columnsString, "c", "", "Columns to replace cells (shorthand)")
	fs.StringVar(&sub.regex, "regex", "", "Regular expression to match for replacement")
	fs.StringVar(&sub.repl, "repl", "", "Replacement string")
	fs.BoolVar(&sub.caseInsensitive, "case-insensitive", false, "Make regex case insensitive")
	fs.BoolVar(&sub.caseInsensitive, "i", false, "Make regex case insensitive (shorthand)")
}

func (sub *ReplaceSubcommand) Run(args []string) {
	// Get columns to compare against
	var columns []string
	if sub.columnsString == "" {
		columns = make([]string, 0)
	} else {
		columns = GetArrayFromCsvString(sub.columnsString)
	}

	// Get replace function
	var replaceFunc func(string) string
	if sub.caseInsensitive {
		sub.regex = "(?i)" + sub.regex
	}
	re, err := regexp.Compile(sub.regex)
	if err != nil {
		ExitWithError(err)
	}
	replaceFunc = func(elem string) string {
		return re.ReplaceAllString(elem, sub.repl)
	}

	inputCsvs := GetInputCsvsOrPanic(args, 1)

	ReplaceWithFunc(inputCsvs[0], columns, replaceFunc)
}

func ReplaceWithFunc(inputCsv *InputCsv, columns []string, replaceFunc func(string) string) {
	outputCsv := NewOutputCsvFromInputCsv(inputCsv)

	// Read header to get column index and write.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}

	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	outputCsv.Write(header)

	// Write replaced rows
	rowToWrite := make([]string, len(header))
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		copy(rowToWrite, row)
		for _, columnIndex := range columnIndices {
			rowToWrite[columnIndex] = replaceFunc(rowToWrite[columnIndex])
		}
		outputCsv.Write(rowToWrite)
	}
}
