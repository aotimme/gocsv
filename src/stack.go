package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"os"
)

type StackSubcommand struct {
	groupName    string
	groupsString string
	useFilenames bool
}

func (sub *StackSubcommand) Name() string {
	return "stack"
}
func (sub *StackSubcommand) Aliases() []string {
	return []string{}
}
func (sub *StackSubcommand) Description() string {
	return "Stack multiple CSVs into one CSV."
}
func (sub *StackSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.groupName, "group-name", "", "Name of the column for grouping")
	fs.StringVar(&sub.groupsString, "groups", "", "Group to display for each file")
	fs.BoolVar(&sub.useFilenames, "filenames", false, "Use the filename for groups")
}

func (sub *StackSubcommand) Run(args []string) {
	filenames := args

	hasSpecifiedGroups := sub.groupsString != ""
	if hasSpecifiedGroups && sub.useFilenames {
		panic(errors.New("Cannot specify both --filename and --groups"))
	}

	shouldAppendGroup := hasSpecifiedGroups || sub.useFilenames

	var groups []string
	if hasSpecifiedGroups {
		groups = GetArrayFromCsvString(sub.groupsString)
	} else if sub.useFilenames {
		groups = filenames
	}

	if shouldAppendGroup && len(filenames) != len(groups) {
		panic(errors.New("Number of files and groups are not equal"))
	}

	var groupColumnName string
	if sub.groupName != "" {
		groupColumnName = sub.groupName
	} else if sub.useFilenames {
		groupColumnName = "File"
	} else if shouldAppendGroup {
		groupColumnName = "Group"
	} else {
		groupColumnName = ""
	}

	inputCsvs := GetInputCsvsOrPanic(filenames, -1)
	StackFiles(inputCsvs, groupColumnName, groups)
}

func StackFiles(inputCsvs []AbstractInputCsv, groupName string, groups []string) {
	shouldAppendGroup := groupName != ""
	writer := csv.NewWriter(os.Stdout)

	// Check that the headers match
	headers := make([][]string, len(inputCsvs))
	for i, inputCsv := range inputCsvs {
		header, err := inputCsv.Read()
		if err != nil {
			panic(err)
		}
		headers[i] = header
	}
	firstHeader := headers[0]
	for i, header := range headers {
		if i == 0 {
			continue
		}
		if len(firstHeader) != len(header) {
			panic(errors.New("Headers do not match"))
		}
		for j, elem := range firstHeader {
			if elem != header[j] {
				panic(errors.New("Headers do not match"))
			}
		}
	}
	if shouldAppendGroup {
		firstHeader = append(firstHeader, groupName)
	}
	writer.Write(firstHeader)
	writer.Flush()

	// Go through the files
	for i, inputCsv := range inputCsvs {
		for {
			row, err := inputCsv.Read()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err)
				}
			}
			if shouldAppendGroup {
				row = append(row, groups[i])
			}
			writer.Write(row)
			writer.Flush()
		}
	}
}
