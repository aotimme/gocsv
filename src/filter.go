package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type FilterSubcommand struct {
	columnsString   string
	exclude         bool
	regex           string
	caseInsensitive bool
	equals          string
	gtStr           string
	gteStr          string
	ltStr           string
	lteStr          string
}

func (sub *FilterSubcommand) Name() string {
	return "filter"
}
func (sub *FilterSubcommand) Aliases() []string {
	return []string{}
}
func (sub *FilterSubcommand) Description() string {
	return "Extract rows whose column match some criterion."
}
func (sub *FilterSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.columnsString, "columns", "", "Columns to filter against")
	fs.StringVar(&sub.columnsString, "c", "", "Columns to filter against (shorthand)")
	fs.BoolVar(&sub.exclude, "exclude", false, "Exclude matching rows")
	fs.StringVar(&sub.regex, "regex", "", "Regular expression for filtering")
	fs.StringVar(&sub.equals, "equals", "", "Exact equality")
	fs.StringVar(&sub.equals, "eq", "", "Exact equality")
	fs.BoolVar(&sub.caseInsensitive, "case-insensitive", false, "Make regular expression case insensitive")
	fs.BoolVar(&sub.caseInsensitive, "i", false, "Make regular expression case insensitive (shorthand)")
	fs.StringVar(&sub.gtStr, "gt", "", "Greater than")
	fs.StringVar(&sub.gteStr, "gte", "", "Greater than or equal to")
	fs.StringVar(&sub.ltStr, "lt", "", "Less than")
	fs.StringVar(&sub.lteStr, "lte", "", "Less than or equal to")
}

func (sub *FilterSubcommand) Run(args []string) {
	// Get columns to compare against
	var columns []string
	if sub.columnsString == "" {
		columns = make([]string, 0)
	} else {
		columns = GetArrayFromCsvString(sub.columnsString)
	}

	// Get match function
	var matchFunc func(string) bool
	if sub.regex != "" {
		if sub.caseInsensitive {
			sub.regex = "(?i)" + sub.regex
		}
		re, err := regexp.Compile(sub.regex)
		if err != nil {
			ExitWithError(err)
		}
		matchFunc = func(elem string) bool {
			return re.MatchString(elem)
		}
	} else if sub.equals != "" {
		matchFunc = func(elem string) bool {
			return elem == sub.equals
		}
	} else if sub.gtStr != "" {
		if IsFloatType(sub.gtStr) {
			gt, err := strconv.ParseFloat(sub.gtStr, 64)
			if err != nil {
				ExitWithError(err)
			}
			matchFunc = func(elem string) bool {
				elem64, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return false
				}
				return elem64 > gt
			}
		} else if IsDateType(sub.gtStr) {
			gt, err := ParseDate(sub.gtStr)
			if err != nil {
				ExitWithError(err)
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
	} else if sub.gteStr != "" {
		if IsFloatType(sub.gteStr) {
			gte, err := strconv.ParseFloat(sub.gteStr, 64)
			if err != nil {
				ExitWithError(err)
			}
			matchFunc = func(elem string) bool {
				elem64, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return false
				}
				return elem64 >= gte
			}
		} else if IsDateType(sub.gteStr) {
			gte, err := ParseDate(sub.gteStr)
			if err != nil {
				ExitWithError(err)
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
	} else if sub.ltStr != "" {
		if IsFloatType(sub.ltStr) {
			lt, err := strconv.ParseFloat(sub.ltStr, 64)
			if err != nil {
				ExitWithError(err)
			}
			matchFunc = func(elem string) bool {
				elem64, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return false
				}
				return elem64 < lt
			}
		} else if IsDateType(sub.ltStr) {
			lt, err := ParseDate(sub.ltStr)
			if err != nil {
				ExitWithError(err)
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
	} else if sub.lteStr != "" {
		if IsFloatType(sub.lteStr) {
			lte, err := strconv.ParseFloat(sub.lteStr, 64)
			if err != nil {
				ExitWithError(err)
			}
			matchFunc = func(elem string) bool {
				elem64, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return false
				}
				return elem64 <= lte
			}
		} else if IsDateType(sub.lteStr) {
			lte, err := ParseDate(sub.lteStr)
			if err != nil {
				ExitWithError(err)
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

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	FilterMatchFunc(inputCsvs[0], columns, sub.exclude, matchFunc)
}

func FilterMatchFunc(inputCsv *InputCsv, columns []string, exclude bool, matchFunc func(string) bool) {
	outputCsv := NewOutputCsvFromInputCsv(inputCsv)

	// Read header to get column index and write.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}

	// Get indices to compare against.
	// If no columns are specified, then check against all.
	columnIndices := GetIndicesForColumnsOrPanic(header, columns)

	outputCsv.Write(header)

	// Write filtered rows.
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
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
			outputCsv.Write(row)
		}
	}
}
