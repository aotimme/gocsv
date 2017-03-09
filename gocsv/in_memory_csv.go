package main

import (
	"encoding/csv"
	"math/rand"
	"time"
)

type InMemoryCsv struct {
	header []string
	rows   [][]string

	// index of column
	isIndexed bool
	index     map[string][]int
}

func NewInMemoryCsv(reader *csv.Reader) *InMemoryCsv {
	imc := new(InMemoryCsv)
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

func (imc *InMemoryCsv) NumRows() int {
	return len(imc.rows)
}

func (imc *InMemoryCsv) NumColumns() int {
	return len(imc.header)
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

func (imc *InMemoryCsv) InferType(columnIndex int) ColumnType {
	curType := INT_TYPE
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

	SortRowsBy(isLessFunc).Sort(imc.rows, reverse)
}

func (imc *InMemoryCsv) SampleRowIndicesWithReplacement(numRows, seed int) []int {
	totalRows := imc.NumRows()
	retval := make([]int, numRows)
	for i := 0; i < numRows; i++ {
		retval[i] = rand.Intn(totalRows)
	}
	return retval
}

func (imc *InMemoryCsv) SampleRowIndicesWithoutReplacement(numRows, seed int) []int {
	totalRows := imc.NumRows()
	permuted := rand.Perm(totalRows)
	retval := make([]int, numRows)
	for i := 0; i < numRows; i++ {
		retval[i] = permuted[i]
	}
	return retval
}

func (imc *InMemoryCsv) SampleRowIndices(numRows int, replace bool, seed int) []int {
	// NOTE: Updating global `rand` variable for the life of the proces...
	if seed != 0 {
		rand.Seed(int64(seed))
	} else {
		rand.Seed(time.Now().UTC().UnixNano())
	}
	if replace {
		return imc.SampleRowIndicesWithReplacement(numRows, seed)
	} else {
		return imc.SampleRowIndicesWithoutReplacement(numRows, seed)
	}
}
