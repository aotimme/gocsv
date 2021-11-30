package cmd

import (
	"fmt"
	"testing"
)

func TestRunAutoincrement(t *testing.T) {
	testCases := []struct {
		name    string
		seed    int
		prepend bool
		rows    [][]string
	}{
		{"ID", 1, false, [][]string{
			{"Number", "String", "ID"},
			{"1", "One", "1"},
			{"2", "Two", "2"},
			{"-1", "Minus One", "3"},
			{"2", "Another Two", "4"},
		}},
		{"ID", 0, false, [][]string{
			{"Number", "String", "ID"},
			{"1", "One", "0"},
			{"2", "Two", "1"},
			{"-1", "Minus One", "2"},
			{"2", "Another Two", "3"},
		}},
		{"PKEY", 1, false, [][]string{
			{"Number", "String", "PKEY"},
			{"1", "One", "1"},
			{"2", "Two", "2"},
			{"-1", "Minus One", "3"},
			{"2", "Another Two", "4"},
		}},
		{"ID", 1, true, [][]string{
			{"ID", "Number", "String"},
			{"1", "1", "One"},
			{"2", "2", "Two"},
			{"3", "-1", "Minus One"},
			{"4", "2", "Another Two"},
		}},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ic, err := NewInputCsv("../test-files/simple-sort.csv")
			if err != nil {
				t.Error("Unexpected error", err)
			}
			toc := new(testOutputCsv)
			sub := new(AutoincrementSubcommand)
			sub.name = tt.name
			sub.seed = tt.seed
			sub.prepend = tt.prepend
			sub.RunAutoincrement(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
