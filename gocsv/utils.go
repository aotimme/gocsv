package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GetIndicesForColumnsOrPanic is a simple wrapper around GetIndicesForColumns
// that will simply panic if GetIndicesForColumns returns an error.
func GetIndicesForColumnsOrPanic(headers []string, columns []string) (indices []int) {
	indices, err := GetIndicesForColumns(headers, columns)
	if err != nil {
		panic(err)
	}
	return
}

// GetIndicesForColumns translates a slice of strings representing the columns requested
// into a slice of the indices of the matching columns.
func GetIndicesForColumns(headers []string, columns []string) (indices []int, err error) {
	if len(columns) == 0 {
		indices = make([]int, len(columns))
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
	if index > -1 {
		panic(fmt.Errorf("Unable to find column specified: %s", column))
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
			return intVal
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
		panic(err)
	}
	return rows[0]
}

// NOTE: Order matters here. Ordered by strictness descending
type ColumnType int

const (
	NULL_TYPE ColumnType = iota
	INT_TYPE
	FLOAT_TYPE
	BOOLEAN_TYPE
	DATE_TYPE
	STRING_TYPE
)

func ColumnTypeToString(columnType ColumnType) string {
	if columnType == NULL_TYPE {
		return "null"
	} else if columnType == INT_TYPE {
		return "int"
	} else if columnType == FLOAT_TYPE {
		return "float"
	} else if columnType == BOOLEAN_TYPE {
		return "boolean"
	} else if columnType == DATE_TYPE {
		return "date"
	} else if columnType == STRING_TYPE {
		return "string"
	} else {
		return ""
	}
}

func ColumnTypeToSqlType(columnType ColumnType) string {
	if columnType == NULL_TYPE {
		return "TEXT"
	} else if columnType == INT_TYPE {
		return "INTEGER"
	} else if columnType == FLOAT_TYPE {
		return "FLOAT"
	} else if columnType == BOOLEAN_TYPE {
		return "INTEGER"
	} else if columnType == DATE_TYPE {
		return "DATE"
	} else if columnType == STRING_TYPE {
		return "TEXT"
	} else {
		return "TEXT"
	}
}

func IsNullType(elem string) bool {
	return elem == ""
}

func IsIntType(elem string) bool {
	_, err := strconv.ParseInt(elem, 0, 0)
	return err == nil
}

func IsFloatType(elem string) bool {
	_, err := strconv.ParseFloat(elem, 64)
	return err == nil
}

func IsBooleanType(elem string) bool {
	strLower := strings.ToLower(elem)
	return strLower == "t" || strLower == "true" || strLower == "f" || strLower == "false"
}

func IsDateType(elem string) bool {
	_, err := ParseDate(elem)
	return err == nil
}

func ParseDate(elem string) (time.Time, error) {
	patterns := []string{
		"2006-01-02",
		"2006-1-2",
		"1/2/2006",
		"01/02/2006",
	}
	for _, pattern := range patterns {
		t, err := time.Parse(pattern, elem)
		if err == nil {
			return t, nil
		}
	}
	return time.Now(), errors.New("Invalid Date string")
}

func ParseDateOrPanic(elem string) time.Time {
	t, err := ParseDate(elem)
	if err != nil {
		panic(err)
	}
	return t
}

func InferType(elem string) ColumnType {
	if IsNullType(elem) {
		return NULL_TYPE
	}
	if IsIntType(elem) {
		return INT_TYPE
	}
	if IsFloatType(elem) {
		return FLOAT_TYPE
	}
	if IsBooleanType(elem) {
		return BOOLEAN_TYPE
	}
	if IsDateType(elem) {
		return DATE_TYPE
	}
	return STRING_TYPE
}

// Adapted from https://golang.org/pkg/sort/#example__sortKeys
type SortRowsBy func(r1, r2 *[]string) bool
type RowSorter struct {
	rows [][]string
	by   func(r1, r2 *[]string) bool
}

func (by SortRowsBy) Sort(rows [][]string, reverse bool) {
	rs := &RowSorter{rows: rows, by: by}
	if reverse {
		sort.Sort(sort.Reverse(rs))
	} else {
		sort.Sort(rs)
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

func ParseFloat64(strVal string) (float64, error) {
	return strconv.ParseFloat(strVal, 64)
}

func ParseFloat64OrPanic(strVal string) float64 {
	floatVal, err := ParseFloat64(strVal)
	if err != nil {
		panic(err)
	}
	return floatVal
}

func ParseInt64(strVal string) (int64, error) {
	return strconv.ParseInt(strVal, 0, 0)
}

func ParseInt64OrPanic(strVal string) int64 {
	intVal, err := ParseInt64(strVal)
	if err != nil {
		panic(err)
	}
	return intVal
}

func GetBaseFilenameWithoutExtension(filename string) string {
	baseFilename := path.Base(filename)
	extension := path.Ext(baseFilename)
	return strings.TrimSuffix(baseFilename, extension)
}
