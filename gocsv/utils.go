package main

import (
  "encoding/csv"
  "errors"
  "fmt"
  "io"
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
  FLOAT_TYPE ColumnType = iota
  INT_TYPE
  STRING_TYPE
)

type InMemoryCsv struct {
  header []string
  rows [][]string

  // index of column
  isIndexed bool
  index map[string][]int
}

func NewInMemoryCsv(r io.Reader) *InMemoryCsv {
  imc := new(InMemoryCsv)
  reader := csv.NewReader(r)
  rows, err := reader.ReadAll()
  if err != nil {
    panic(err)
  }
  imc.header = rows[0]
  imc.rows = rows[1:]
  imc.isIndexed = false
  return imc
}

func (imc *InMemoryCsv) Index(columnIndex int) {
  imc.index = make(map[string][]int)
  for i, row := range imc.rows {
    rowval := row[columnIndex]
    group, ok := imc.index[rowval]
    if ok {
      group = append(group, i)
    } else {
      group = make([]int, 1)
      group[0] = i
    }
    imc.index[rowval] = group
  }
}


func (imc *InMemoryCsv) GetRowIndicesMatchingIndexedColumn(value string) []int {
  indices, ok := imc.index[value]
  if ok {
    return indices
  } else {
    return make([]int, 0)
  }
}

func (imc *InMemoryCsv) GetRowsMatchingIndexedColumn(value string) [][]string {
  indices := imc.GetRowIndicesMatchingIndexedColumn(value)
  rows := make([][]string, 0)
  for _, idx := range indices {
    rows = append(rows, imc.rows[idx])
  }
  return rows
}

func inferType(elem string) ColumnType {
  _, err := strconv.ParseFloat(elem, 64)
  if err == nil {
    return FLOAT_TYPE
  }
  err = nil
  _, err = strconv.ParseInt(elem, 0, 0)
  if err == nil {
    return INT_TYPE
  }
  return STRING_TYPE
}

func (imc *InMemoryCsv) InferType(columnIndex int) ColumnType {
  curType := FLOAT_TYPE
  for _, row := range imc.rows {
    thisType := inferType(row[columnIndex])
    if thisType > curType {
      curType = thisType
    }
    if curType == STRING_TYPE {
      return curType
    }
  }
  return curType
}

// Adapted from https://golang.org/pkg/sort/#example__sortKeys
type SortRowsBy func(r1, r2 *[]string) bool
type RowSorter struct {
  rows [][]string
  by func(r1, r2 *[]string) bool
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

func (imc *InMemoryCsv) SortRows(columnIndices []int, columnTypes []ColumnType, reverse bool) {
  isLessFunc := func(row1Ptr, row2Ptr *[]string) bool {
    row1 := *row1Ptr
    row2 := *row2Ptr
    for i, columnIndex := range columnIndices {
      columnType := columnTypes[i]
      if columnType == FLOAT_TYPE {
        row1Val := parseFloat64OrPanic(row1[columnIndex])
        row2Val := parseFloat64OrPanic(row2[columnIndex])
        if row1Val < row2Val {
          return true
        } else if row1Val > row2Val {
          return false
        }
      } else if columnType == INT_TYPE {
        row1Val := parseInt64OrPanic(row1[columnIndex])
        row2Val := parseInt64OrPanic(row2[columnIndex])
        if row1Val < row2Val {
          return true
        } else if row1Val > row2Val {
          return false
        }
      } else {
        row1Val := row1[columnIndex]
        row2Val := row2[columnIndex]
        if row1Val < row2Val {
          return true
        } else if row1Val > row2Val {
          return false
        }
      }
    }
    return true
  }

  // TODO: Respect "reverse"
  SortRowsBy(isLessFunc).Sort(imc.rows, reverse)
}
