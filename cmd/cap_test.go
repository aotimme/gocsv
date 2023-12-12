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
		{"", false, "Col", [][]string{
			{"Col", "Col 1"},
			{"Number", "String"},
			{"1", "One"},
			{"2", "Two"},
			{"-1", "Minus One"},
			{"2", "Another Two"},
		}},
		{"Numero,Cadena", false, "", [][]string{
			{"Numero", "Cadena"},
			{"Number", "String"},
			{"1", "One"},
			{"2", "Two"},
			{"-1", "Minus One"},
			{"2", "Another Two"},
		}},
		{"Numero,Cadena,Otro", true, "", [][]string{
			{"Numero", "Cadena"},
			{"Number", "String"},
			{"1", "One"},
			{"2", "Two"},
			{"-1", "Minus One"},
			{"2", "Another Two"},
		}},
		{"Numero", false, "Column", [][]string{
			{"Numero", "Column"},
			{"Number", "String"},
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
