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
			[]string{"Number", "String", "ID"},
			[]string{"1", "One", "1"},
			[]string{"2", "Two", "2"},
			[]string{"-1", "Minus One", "3"},
			[]string{"2", "Another Two", "4"},
		}},
		{"ID", 0, false, [][]string{
			[]string{"Number", "String", "ID"},
			[]string{"1", "One", "0"},
			[]string{"2", "Two", "1"},
			[]string{"-1", "Minus One", "2"},
			[]string{"2", "Another Two", "3"},
		}},
		{"PKEY", 1, false, [][]string{
			[]string{"Number", "String", "PKEY"},
			[]string{"1", "One", "1"},
			[]string{"2", "Two", "2"},
			[]string{"-1", "Minus One", "3"},
			[]string{"2", "Another Two", "4"},
		}},
		{"ID", 1, true, [][]string{
			[]string{"ID", "Number", "String"},
			[]string{"1", "1", "One"},
			[]string{"2", "2", "Two"},
			[]string{"3", "-1", "Minus One"},
			[]string{"4", "2", "Another Two"},
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
