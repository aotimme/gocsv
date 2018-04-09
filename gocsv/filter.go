package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type FilterSubcommand struct{}

func (sub *FilterSubcommand) Name() string {
	return "filter"
}
func (sub *FilterSubcommand) Aliases() []string {
	return []string{}
}
func (sub *FilterSubcommand) Description() string {
	return "Extract rows whose column match some criterion."
}

func (sub *FilterSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	var regex, columnsString string
	var exclude, caseInsensitive bool
	var gtStr, gteStr, ltStr, lteStr string
	fs.StringVar(&columnsString, "columns", "", "Columns to filter against")
	fs.StringVar(&columnsString, "c", "", "Columns to filter against (shorthand)")
	fs.BoolVar(&exclude, "exclude", false, "Exclude matching rows")
	fs.StringVar(&regex, "regex", "", "Regular expression for filtering")
	fs.BoolVar(&caseInsensitive, "case-insensitive", false, "Make regular expression case insensitive")
	fs.BoolVar(&caseInsensitive, "i", false, "Make regular expression case insensitive (shorthand)")
	fs.StringVar(&gtStr, "gt", "", "Greater than")
	fs.StringVar(&gteStr, "gte", "", "Greater than or equal to")
	fs.StringVar(&ltStr, "lt", "", "Less than")
	fs.StringVar(&lteStr, "lte", "", "Less than or equal to")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	// Get columns to compare against
	var columns []string
	if columnsString == "" {
		columns = make([]string, 0)
	} else {
		columns = GetArrayFromCsvString(columnsString)
	}

	// Get match function
	var matchFunc func(string) bool
	if regex != "" {
		if caseInsensitive {
			regex = "(?i)" + regex
		}
		re, err := regexp.Compile(regex)
		if err != nil {
			panic(err)
		}
		matchFunc = func(elem string) bool {
			return re.MatchString(elem)
		}
	} else if gtStr != "" {
		if IsFloatType(gtStr) {
			gt, err := strconv.ParseFloat(gtStr, 64)
			if err != nil {
				panic(err)
			}
			matchFunc = func(elem string) bool {
				elem64, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return false
				}
				return elem64 > gt
			}
		} else if IsDateType(gtStr) {
			gt, err := ParseDate(gtStr)
			if err != nil {
				panic(err)
			}
			matchFunc = func(elem string) bool {
				elemDate, err := ParseDate(elem)
				if err != nil {
					return false
				}
				return elemDate.After(gt)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Invalid argument for -gt")
			os.Exit(1)
		}
	} else if gteStr != "" {
		if IsFloatType(gteStr) {
			gte, err := strconv.ParseFloat(gteStr, 64)
			if err != nil {
				panic(err)
			}
			matchFunc = func(elem string) bool {
				elem64, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return false
				}
				return elem64 >= gte
			}
		} else if IsDateType(gteStr) {
			gte, err := ParseDate(gteStr)
			if err != nil {
				panic(err)
			}
			matchFunc = func(elem string) bool {
				elemDate, err := ParseDate(elem)
				if err != nil {
					return false
				}
				return elemDate.Equal(gte) || elemDate.After(gte)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Invalid argument for -gte")
			os.Exit(1)
		}
	} else if ltStr != "" {
		if IsFloatType(ltStr) {
			lt, err := strconv.ParseFloat(ltStr, 64)
			if err != nil {
				panic(err)
			}
			matchFunc = func(elem string) bool {
				elem64, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return false
				}
				return elem64 < lt
			}
		} else if IsDateType(ltStr) {
			lt, err := ParseDate(ltStr)
			if err != nil {
				panic(err)
			}
			matchFunc = func(elem string) bool {
				elemDate, err := ParseDate(elem)
				if err != nil {
					return false
				}
				return elemDate.Before(lt)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Invalid argument for -lt")
			os.Exit(1)
		}
	} else if lteStr != "" {
		if IsFloatType(lteStr) {
			lte, err := strconv.ParseFloat(lteStr, 64)
			if err != nil {
				panic(err)
			}
			matchFunc = func(elem string) bool {
				elem64, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return false
				}
				return elem64 <= lte
			}
		} else if IsDateType(lteStr) {
			lte, err := ParseDate(lteStr)
			if err != nil {
				panic(err)
			}
			matchFunc = func(elem string) bool {
				elemDate, err := ParseDate(elem)
				if err != nil {
					return false
				}
				return elemDate.Equal(lte) || elemDate.Before(lte)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Invalid argument for -lte")
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Missing filter function")
		os.Exit(1)
	}

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	FilterMatchFunc(inputCsvs[0], columns, exclude, matchFunc)
}

func FilterMatchFunc(inputCsv AbstractInputCsv, columns []string, exclude bool, matchFunc func(string) bool) {
	writer := csv.NewWriter(os.Stdout)

	// Read header to get column index and write.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}

	// Get indices to compare against.
	// If no columns are specified, then check against all.
	var columnIndices []int
	if len(columns) == 0 {
		columnIndices = make([]int, len(header))
		for i, _ := range header {
			columnIndices[i] = i
		}
	} else {
		columnIndices = make([]int, len(columns))
		for i, column := range columns {
			index := GetColumnIndexOrPanic(header, column)
			columnIndices[i] = index
		}
	}

	writer.Write(header)
	writer.Flush()

	// Write filtered rows.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		rowMatches := false
		for _, columnIndex := range columnIndices {
			if matchFunc(row[columnIndex]) {
				rowMatches = true
				break
			}
		}
		shouldOutputRow := (!exclude && rowMatches) || (exclude && !rowMatches)
		if shouldOutputRow {
			writer.Write(row)
			writer.Flush()
		}
	}
}
