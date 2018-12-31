package cmd

import (
	"bytes"
	"flag"
	"io"
	"strconv"
	"text/template"
)

type AddSubcommand struct {
	name     string
	template string
	prepend  bool
}

func (sub *AddSubcommand) Name() string {
	return "add"
}
func (sub *AddSubcommand) Aliases() []string {
	// Adding "template" and "tmpl" for backwards compatibility
	return []string{"template", "tmpl"}
}
func (sub *AddSubcommand) Description() string {
	return "Add a column to a CSV."
}
func (sub *AddSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.name, "name", "", "Name of new column")
	fs.StringVar(&sub.name, "n", "", "Name of new column (shorthand)")
	fs.StringVar(&sub.template, "template", "", "Template for the new column")
	fs.StringVar(&sub.template, "t", "", "Template for the new column (shorthand)")
	fs.BoolVar(&sub.prepend, "prepend", false, "Prepend the new column (defaults to append)")
}

func (sub *AddSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsv(inputCsvs[0])
	sub.RunAdd(inputCsvs[0], outputCsv)
}

func (sub *AddSubcommand) RunAdd(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	tmpl, err := template.New("template").Parse(sub.template)
	if err != nil {
		ExitWithError(err)
	}
	AddColumn(inputCsv, outputCsvWriter, tmpl, sub.name, sub.prepend)
	err = inputCsv.Close()
	if err != nil {
		ExitWithError(err)
	}
}

func AddColumn(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, tmpl *template.Template, name string, prepend bool) {
	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}

	numInputColumns := len(header)
	shellRow := make([]string, numInputColumns+1)
	if prepend {
		shellRow[0] = name
		for i, elem := range header {
			shellRow[i+1] = elem
		}
	} else {
		copy(shellRow, header)
		shellRow[numInputColumns] = name
	}
	outputCsvWriter.Write(shellRow)

	// Create the holding map for the template data.
	templateData := make(map[string]string)

	// Write rows with template.
	index := 1
	for {
		templateData["index"] = strconv.Itoa(index)
		index++

		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		for i, elem := range row {
			templateData[header[i]] = elem
		}

		var rendered bytes.Buffer
		err = tmpl.Execute(&rendered, templateData)
		if err != nil {
			ExitWithError(err)
		}

		newElem := rendered.String()

		if prepend {
			shellRow[0] = newElem
			for i, elem := range row {
				shellRow[i+1] = elem
			}
		} else {
			copy(shellRow, row)
			shellRow[numInputColumns] = newElem
		}
		outputCsvWriter.Write(shellRow)
	}
}
