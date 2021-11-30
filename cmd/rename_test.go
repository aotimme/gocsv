package cmd

import (
	"fmt"
	"testing"
)

func TestRunRename(t *testing.T) {
	testCases := []struct {
		columnsString string
		namesString   string
		rows          [][]string
	}{
		{"Number,String", "Numero,Cadena", [][]string{
			{"Numero", "Cadena"},
			{"1", "One"},
			{"2", "Two"},
			{"-1", "Minus One"},
			{"2", "Another Two"},
		}},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ic, err := NewInputCsv("../test-files/simple-sort.csv")
			if err != nil {
				t.Error("Unexpected error", err)
			}
			toc := new(testOutputCsv)
			sub := new(RenameSubcommand)
			sub.columnsString = tt.columnsString
			sub.namesString = tt.namesString
			sub.RunRename(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
