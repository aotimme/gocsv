package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"os"
)

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

func RunStack(args []string) {
	fs := flag.NewFlagSet("stack", flag.ExitOnError)
	var groupName, groupsString string
	var useFilenames bool
	fs.StringVar(&groupName, "group-name", "", "Name of the column for grouping")
	fs.StringVar(&groupsString, "groups", "", "Group to display for each file")
	fs.BoolVar(&useFilenames, "filenames", false, "Use the filename for groups")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	filenames := fs.Args()

	hasSpecifiedGroups := groupsString != ""
	if hasSpecifiedGroups && useFilenames {
		panic(errors.New("Cannot specify both --filename and --groups"))
	}

	shouldAppendGroup := hasSpecifiedGroups || useFilenames

	var groups []string
	if hasSpecifiedGroups {
		groups = GetArrayFromCsvString(groupsString)
	} else if useFilenames {
		groups = filenames
	}

	if shouldAppendGroup && len(filenames) != len(groups) {
		panic(errors.New("Number of files and groups are not equal"))
	}

	var groupColumnName string
	if groupName != "" {
		groupColumnName = groupName
	} else if useFilenames {
		groupColumnName = "File"
	} else if shouldAppendGroup {
		groupColumnName = "Group"
	} else {
		groupColumnName = ""
	}

	inputCsvs, err := GetInputCsvs(filenames, -1)
	if err != nil {
		panic(err)
	}
	StackFiles(inputCsvs, groupColumnName, groups)
}
