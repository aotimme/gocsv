package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/template"
)

type TemplateSubcommand struct {
	name     string
	template string
	prepend  bool
}

func (sub *TemplateSubcommand) Name() string {
	return "template"
}
func (sub *TemplateSubcommand) Aliases() []string {
	return []string{"tmpl"}
}
func (sub *TemplateSubcommand) Description() string {
	return "Add a column with values based on a template using other columns."
}
func (sub *TemplateSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.name, "name", "Templated", "Name of templated column")
	fs.StringVar(&sub.template, "template", "", "Template for column")
	fs.StringVar(&sub.template, "t", "", "Template for column (shorthand)")
	fs.BoolVar(&sub.prepend, "prepend", false, "Prepend the templated column (defaults to append)")
}

func (sub *TemplateSubcommand) Run(args []string) {
	if sub.template == "" {
		fmt.Fprintln(os.Stderr, "Missing argument --template (-t)")
		os.Exit(1)
	}
	tmpl, err := template.New("template").Parse(sub.template)
	if err != nil {
		ExitWithError(err)
	}
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	RenderTemplate(inputCsvs[0], tmpl, sub.name, sub.prepend)
	err = inputCsvs[0].Close()
	if err != nil {
		ExitWithError(err)
	}
}

func RenderTemplate(inputCsv *InputCsv, tmpl *template.Template, name string, prepend bool) {
	outputCsv := NewOutputCsvFromInputCsv(inputCsv)

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
	outputCsv.Write(shellRow)

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
		outputCsv.Write(shellRow)
	}
}
