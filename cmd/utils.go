package cmd

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/aotimme/gocsv/csv"
)

const (
	BOM_RUNE      = '\uFEFF'
	BOM_STRING    = "\uFEFF"
	NUM_BOM_BYTES = 3
)

func GetDelimiterFromString(delimiter string) rune {
	if delimiter == "\\t" {
		return '\t'
	} else if len(delimiter) > 0 {
		delimiterRune, _ := utf8.DecodeRuneInString(delimiter)
		return delimiterRune
	}
	return rune(0)
}

// GetIndicesForColumnsOrPanic is a simple wrapper around GetIndicesForColumns
// that will simply panic if GetIndicesForColumns returns an error.
func GetIndicesForColumnsOrPanic(headers []string, columns []string) (indices []int) {
	indices, err := GetIndicesForColumns(headers, columns)
	if err != nil {
		ExitWithError(err)
	}
	return
}

// GetIndicesForColumns translates a slice of strings representing the columns requested
// into a slice of the indices of the matching columns.
func GetIndicesForColumns(headers []string, columns []string) (indices []int, err error) {
	if len(columns) == 0 {
		indices = make([]int, len(headers))
		for i := range indices {
			indices[i] = i
		}
		return
	}
	for _, column := range columns {
		columnIndices, err := GetIndicesForColumn(headers, column)
		if err != nil {
			return nil, err
		}
		for _, columnIndex := range columnIndices {
			indices = append(indices, columnIndex)
		}
	}
	return
}

// GetIndicesForColumn translates a string representing a column requested
// into a slice of the indices of the matching column.
func GetIndicesForColumn(headers []string, column string) (indices []int, err error) {
	int64Val, errParse := strconv.ParseInt(column, 0, 0)
	if errParse == nil {
		intVal := int(int64Val)
		if intVal > 0 && intVal <= len(headers) {
			return []int{intVal - 1}, nil
		}
	}
	possibleIntStrs := strings.Split(column, "-")
	if len(possibleIntStrs) == 2 {
		minVal64, err1 := parsePossibleIntInHeader(possibleIntStrs[0], 1)
		maxVal64, err2 := parsePossibleIntInHeader(possibleIntStrs[1], int64(len(headers)))
		if err1 == nil && err2 == nil {
			minVal := int(minVal64)
			maxVal := int(maxVal64)
			ascending := true
			if minVal > maxVal {
				ascending = false
				minVal, maxVal = maxVal, minVal
			}
			if minVal <= maxVal && minVal > 0 && minVal <= len(headers) && maxVal <= len(headers) {
				indices = make([]int, maxVal-minVal+1)
				for i := range indices {
					if ascending {
						indices[i] = minVal + i - 1
					} else {
						indices[i] = maxVal - i - 1
					}
				}
				return
			}
		}
	}
	indices = GetIndicesOfString(headers, column)
	if len(indices) == 0 {
		err = fmt.Errorf("Could not find header \"%s\"", column)
		return
	}
	return
}

func parsePossibleIntInHeader(possibleIntStr string, valueIfEmpty int64) (int64, error) {
	if possibleIntStr == "" {
		return valueIfEmpty, nil
	}
	return strconv.ParseInt(possibleIntStr, 0, 0)
}

// GetIndicesOfString returns a slice of the indices at which the passed in string
// occurs in the input slice.
func GetIndicesOfString(haystack []string, needle string) (indices []int) {
	for i, cand := range haystack {
		if cand == needle {
			indices = append(indices, i)
		}
	}
	return
}

// GetIndexForColumnOrPanic is a simple wrapper around GetIndexForColumn
// that will simply panic if GetIndexForColumn returns -1.
func GetIndexForColumnOrPanic(headers []string, column string) int {
	index := GetIndexForColumn(headers, column)
	if index == -1 {
		ExitWithError(fmt.Errorf("Unable to find column specified: %s", column))
	}
	return index
}

// GetIndexForColumn finds the single index of a header given a column specification
// Note that this method assumes that only one index is requested so it has slightly
// different logic from GetIndicesForColumn.
func GetIndexForColumn(headers []string, column string) int {
	int64Val, errParse := strconv.ParseInt(column, 0, 0)
	if errParse == nil {
		intVal := int(int64Val)
		if intVal > 0 && intVal <= len(headers) {
			return intVal - 1
		}
	}
	return GetFirstIndexOfString(headers, column)
}

// GetFirstIndexOfString gets the index of the first occurrence of
// a string in a string slice.
func GetFirstIndexOfString(haystack []string, needle string) int {
	for i, cand := range haystack {
		if cand == needle {
			return i
		}
	}
	return -1
}

func GetArrayFromCsvString(s string) []string {
	c := csv.NewReader(strings.NewReader(s))
	rows, err := c.ReadAll()
	if err != nil {
		ExitWithError(err)
	}
	return rows[0]
}

// Adapted from https://golang.org/pkg/sort/#example__sortKeys
type SortRowsBy func(r1, r2 *[]string) bool
type RowSorter struct {
	rows [][]string
	by   func(r1, r2 *[]string) bool
}

func (by SortRowsBy) Sort(rows [][]string, stable bool, reverse bool) {
	rs := &RowSorter{rows: rows, by: by}
	var sortFunc = sort.Sort
	if stable {
		sortFunc = sort.Stable
	}
	if reverse {
		sortFunc(sort.Reverse(rs))
	} else {
		sortFunc(rs)
	}
}
func (rs *RowSorter) Len() int {
	return len(rs.rows)
}
func (rs *RowSorter) Swap(i, j int) {
	rs.rows[i], rs.rows[j] = rs.rows[j], rs.rows[i]
}
func (rs *RowSorter) Less(i, j int) bool {
	return rs.by(&rs.rows[i], &rs.rows[j])
}

func GetBaseFilenameWithoutExtension(filename string) string {
	baseFilename := path.Base(filename)
	extension := path.Ext(baseFilename)
	return strings.TrimSuffix(baseFilename, extension)
}

func concat(outrow, row1, row2 []string) {
	i := 0
	for _, elem := range row1 {
		outrow[i] = elem
		i++
	}
	for _, elem := range row2 {
		outrow[i] = elem
		i++
	}
}

func ExitWithError(err error) {
	if DEBUG {
		panic(err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
