package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"./csv"
)

type JoinSubcommand struct {
	columnsString string
	left          bool
	right         bool
	outer         bool
}

func (sub *JoinSubcommand) Name() string {
	return "join"
}
func (sub *JoinSubcommand) Aliases() []string {
	return []string{}
}
func (sub *JoinSubcommand) Description() string {
	return "Join two CSVs based on equality of elements in a column."
}
func (sub *JoinSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.columnsString, "columns", "", "Columns to join on")
	fs.StringVar(&sub.columnsString, "c", "", "Columns to join on (shorthand)")
	fs.BoolVar(&sub.left, "left", false, "Left join")
	fs.BoolVar(&sub.right, "right", false, "Right join")
	fs.BoolVar(&sub.outer, "outer", false, "Full outer join")
}

func (sub *JoinSubcommand) Run(args []string) {
	if sub.columnsString == "" {
		fmt.Fprintln(os.Stderr, "Missing required argument --columns")
		os.Exit(1)
	}
	numJoins := 0
	if sub.left {
		numJoins++
	}
	if sub.right {
		numJoins++
	}
	if sub.outer {
		numJoins++
	}
	if numJoins > 1 {
		fmt.Fprintln(os.Stderr, "Must only specify zero or one of --left, --right, or --outer")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(sub.columnsString)
	if len(columns) < 1 || len(columns) > 2 {
		fmt.Fprintln(os.Stderr, "Invalid argument for --columns")
		os.Exit(1)
	}
	if len(columns) == 1 {
		columns = append(columns, columns[0])
	}

	inputCsvs := GetInputCsvsOrPanic(args, 2)

	if sub.left {
		LeftJoin(inputCsvs[0], inputCsvs[1], columns[0], columns[1])
	} else if sub.right {
		RightJoin(inputCsvs[0], inputCsvs[1], columns[0], columns[1])
	} else if sub.outer {
		OuterJoin(inputCsvs[0], inputCsvs[1], columns[0], columns[1])
	} else {
		InnerJoin(inputCsvs[0], inputCsvs[1], columns[0], columns[1])
	}
}

func InnerJoin(leftInputCsv, rightInputCsv AbstractInputCsv, leftColname, rightColname string) {
	leftHeader, err := leftInputCsv.Read()
	if err != nil {
		panic(err)
	}
	leftColIndex := GetIndexForColumnOrPanic(leftHeader, leftColname)
	numLeftColumns := len(leftHeader)

	rightCsv := NewInMemoryCsvFromInputCsv(rightInputCsv)
	rightColIndex := GetIndexForColumnOrPanic(rightCsv.header, rightColname)
	numRightColumns := len(rightCsv.header)
	rightCsv.Index(rightColIndex)

	shellRow := make([]string, numLeftColumns+numRightColumns)

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	concat(shellRow, leftHeader, rightCsv.header)
	writer.Write(shellRow)
	writer.Flush()

	// Write inner-joined rows.
	for {
		row, err := leftInputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		rightRows := rightCsv.GetRowsMatchingIndexedColumn(row[leftColIndex])
		if len(rightRows) > 0 {
			for _, rightRow := range rightRows {
				concat(shellRow, row, rightRow)
				writer.Write(shellRow)
				writer.Flush()
			}
		}
	}
}

func LeftJoin(leftInputCsv, rightInputCsv AbstractInputCsv, leftColname, rightColname string) {
	leftHeader, err := leftInputCsv.Read()
	if err != nil {
		panic(err)
	}
	leftColIndex := GetIndexForColumnOrPanic(leftHeader, leftColname)
	numLeftColumns := len(leftHeader)

	rightCsv := NewInMemoryCsvFromInputCsv(rightInputCsv)
	rightColIndex := GetIndexForColumnOrPanic(rightCsv.header, rightColname)
	numRightColumns := len(rightCsv.header)
	rightCsv.Index(rightColIndex)

	emptyRightRow := make([]string, numRightColumns)
	shellRow := make([]string, numLeftColumns+numRightColumns)

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	concat(shellRow, leftHeader, rightCsv.header)
	writer.Write(shellRow)
	writer.Flush()

	// Write left-joined rows.
	for {
		row, err := leftInputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		rightRows := rightCsv.GetRowsMatchingIndexedColumn(row[leftColIndex])
		if len(rightRows) > 0 {
			for _, rightRow := range rightRows {
				concat(shellRow, row, rightRow)
				writer.Write(shellRow)
				writer.Flush()
			}
		} else {
			concat(shellRow, row, emptyRightRow)
			writer.Write(shellRow)
			writer.Flush()
		}
	}
}

func RightJoin(leftInputCsv, rightInputCsv AbstractInputCsv, leftColname, rightColname string) {
	rightHeader, err := rightInputCsv.Read()
	if err != nil {
		panic(err)
	}
	rightColIndex := GetIndexForColumnOrPanic(rightHeader, rightColname)
	numRightColumns := len(rightHeader)

	leftCsv := NewInMemoryCsvFromInputCsv(leftInputCsv)
	leftColIndex := GetIndexForColumnOrPanic(leftCsv.header, leftColname)
	leftCsv.Index(leftColIndex)
	numLeftColumns := len(leftCsv.header)

	emptyLeftRow := make([]string, numLeftColumns)
	shellRow := make([]string, numLeftColumns+numRightColumns)

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	concat(shellRow, leftCsv.header, rightHeader)
	writer.Write(shellRow)
	writer.Flush()

	// Write right-joined rows.
	for {
		row, err := rightInputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		leftRows := leftCsv.GetRowsMatchingIndexedColumn(row[rightColIndex])
		if len(leftRows) > 0 {
			for _, leftRow := range leftRows {
				concat(shellRow, leftRow, row)
				writer.Write(shellRow)
				writer.Flush()
			}
		} else {
			concat(shellRow, emptyLeftRow, row)
			writer.Write(shellRow)
			writer.Flush()
		}
	}
}

func OuterJoin(leftInputCsv, rightInputCsv AbstractInputCsv, leftColname, rightColname string) {
	// Basically do a left join and then append any rows from the right table
	// that weren't already included.

	leftHeader, err := leftInputCsv.Read()
	if err != nil {
		panic(err)
	}
	leftColIndex := GetIndexForColumnOrPanic(leftHeader, leftColname)
	numLeftColumns := len(leftHeader)

	rightCsv := NewInMemoryCsvFromInputCsv(rightInputCsv)
	rightColIndex := GetIndexForColumnOrPanic(rightCsv.header, rightColname)
	numRightColumns := len(rightCsv.header)
	rightCsv.Index(rightColIndex)

	emptyLeftRow := make([]string, numLeftColumns)
	emptyRightRow := make([]string, numRightColumns)
	shellRow := make([]string, numLeftColumns+numRightColumns)

	// whether the row in the right column has been included already.
	rightIncludeStatus := make([]bool, len(rightCsv.rows))

	writer := csv.NewWriter(os.Stdout)

	// Write header.
	concat(shellRow, leftHeader, rightCsv.header)
	writer.Write(shellRow)
	writer.Flush()

	// Write left-joined rows.
	for {
		row, err := leftInputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		rightRowIndices := rightCsv.GetRowIndicesMatchingIndexedColumn(row[leftColIndex])
		if len(rightRowIndices) > 0 {
			for _, rightRowIndex := range rightRowIndices {
				rightIncludeStatus[rightRowIndex] = true
				concat(shellRow, row, rightCsv.rows[rightRowIndex])
				writer.Write(shellRow)
				writer.Flush()
			}
		} else {
			concat(shellRow, row, emptyRightRow)
			writer.Write(shellRow)
			writer.Flush()
		}
	}

	// Write remaining right rows.
	for i, row := range rightCsv.rows {
		if rightIncludeStatus[i] {
			continue
		}
		concat(shellRow, emptyLeftRow, row)
		writer.Write(shellRow)
		writer.Flush()
	}
}
