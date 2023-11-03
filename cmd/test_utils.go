package cmd

import "fmt"

type testOutputCsv struct {
	rows [][]string
}

func (toc *testOutputCsv) Write(row []string) error {
	newRow := make([]string, len(row))
	copy(newRow, row)
	toc.rows = append(toc.rows, newRow)
	return nil
}

func assertRowsEqual(expectedRows, actualRows [][]string) error {
	if len(expectedRows) != len(actualRows) {
		return fmt.Errorf("expected %d rows but got %d", len(expectedRows), len(actualRows))
	}
	for i, row := range expectedRows {
		if len(row) != len(actualRows[i]) {
			return fmt.Errorf("expected row %d to have %d entries but got %d", i, len(row), len(actualRows[i]))
		}
		for j, cell := range row {
			if cell != actualRows[i][j] {
				return fmt.Errorf("expected %s in cell (%d, %d) but got %s", cell, i, j, actualRows[i][j])
			}
		}
	}
	return nil
}
