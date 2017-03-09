package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
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
		panic(errors.New(fmt.Sprintf("No column matching \"%s\"", columnIndex)))
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
	INT_TYPE ColumnType = iota
	FLOAT_TYPE
	BOOLEAN_TYPE
	STRING_TYPE
)

func ColumnTypeToString(columnType ColumnType) string {
	if columnType == INT_TYPE {
		return "int"
	} else if columnType == FLOAT_TYPE {
		return "float"
	} else if columnType == BOOLEAN_TYPE {
		return "boolean"
	} else if columnType == STRING_TYPE {
		return "string"
	} else {
		return ""
	}
}

func inferType(elem string) ColumnType {
	_, err := strconv.ParseInt(elem, 0, 0)
	if err == nil {
		return INT_TYPE
	}
	err = nil
	_, err = strconv.ParseFloat(elem, 64)
	if err == nil {
		return FLOAT_TYPE
	}
	strLower := strings.ToLower(strings.Trim(elem, " "))
	if strLower == "t" || strLower == "true" || strLower == "f" || strLower == "false" {
		return BOOLEAN_TYPE
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

func parseFloat64OrPanic(strVal string) float64 {
	floatVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		panic(err)
	}
	return floatVal
}

func parseInt64OrPanic(strVal string) int64 {
	intVal, err := strconv.ParseInt(strVal, 0, 0)
	if err != nil {
		panic(err)
	}
	return intVal
}
