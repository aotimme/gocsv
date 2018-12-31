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
			[]string{"Numero", "Cadena"},
			[]string{"1", "One"},
			[]string{"2", "Two"},
			[]string{"-1", "Minus One"},
			[]string{"2", "Another Two"},
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
