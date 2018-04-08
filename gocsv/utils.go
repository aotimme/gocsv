package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetIndexOfString(haystack []string, needle string) int {
	for i, cand := range haystack {
		if cand == needle {
			return i
		}
	}
	return -1
}

func GetColumnIndexOrPanic(headers []string, header string) int {
	columnIndex := GetIndexOfStringOrIndexMinusOne(headers, header)
	if columnIndex == -1 {
		panic(errors.New(fmt.Sprintf("No column matching \"%d\"", columnIndex)))
	}
	return columnIndex
}

func GetIndexOfStringOrIndexMinusOne(haystack []string, needle string) int {
	int64Val, err := strconv.ParseInt(needle, 0, 0)
	if err == nil {
		intVal := int(int64Val)
		if intVal > 0 && intVal <= len(haystack) {
			return intVal - 1
		}
	}
	return GetIndexOfString(haystack, needle)
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
