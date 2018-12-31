package cmd

import (
	"fmt"
	"testing"
)

func TestRunCap(t *testing.T) {
	testCases := []struct {
		namesString   string
		truncateNames bool
		defaultName   string
		rows          [][]string
	}{
		{"Numero,Cadena", false, "", [][]string{
			[]string{"Numero", "Cadena"},
			[]string{"Number", "String"},
			[]string{"1", "One"},
			[]string{"2", "Two"},
			[]string{"-1", "Minus One"},
			[]string{"2", "Another Two"},
		}},
		{"Numero,Cadena,Otro", true, "", [][]string{
			[]string{"Numero", "Cadena"},
			[]string{"Number", "String"},
			[]string{"1", "One"},
			[]string{"2", "Two"},
			[]string{"-1", "Minus One"},
			[]string{"2", "Another Two"},
		}},
		{"Numero", false, "Column", [][]string{
			[]string{"Numero", "Column"},
			[]string{"Number", "String"},
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
			sub := new(CapSubcommand)
			sub.namesString = tt.namesString
			sub.truncateNames = tt.truncateNames
			sub.defaultName = tt.defaultName
			sub.RunCap(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
