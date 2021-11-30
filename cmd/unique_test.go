package cmd

import (
	"fmt"
	"testing"
)

func TestRunUnique(t *testing.T) {
	testCases := []struct {
		columnsString string
		sorted        bool
		count         bool
		rows          [][]string
	}{
		{"Number", false, false, [][]string{
			{"Number", "String"},
			{"1", "One"},
			{"2", "Two"},
			{"-1", "Minus One"},
		}},
		{"Number,String", false, false, [][]string{
			{"Number", "String"},
			{"1", "One"},
			{"2", "Two"},
			{"-1", "Minus One"},
			{"2", "Another Two"},
		}},
		{"Number", true, false, [][]string{
			{"Number", "String"},
			{"1", "One"},
			{"2", "Two"},
			{"-1", "Minus One"},
			{"2", "Another Two"},
		}},
		{"Number", false, true, [][]string{
			{"Number", "String", "Count"},
			{"1", "One", "1"},
			{"2", "Two", "2"},
			{"-1", "Minus One", "1"},
		}},
		// not actually sorted, so won't notice the duplicates
		{"Number", true, true, [][]string{
			{"Number", "String", "Count"},
			{"1", "One", "1"},
			{"2", "Two", "1"},
			{"-1", "Minus One", "1"},
			{"2", "Another Two", "1"},
		}},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ic, err := NewInputCsv("../test-files/simple-sort.csv")
			if err != nil {
				t.Error("Unexpected error", err)
			}
			toc := new(testOutputCsv)
			sub := new(UniqueSubcommand)
			sub.columnsString = tt.columnsString
			sub.sorted = tt.sorted
			sub.count = tt.count
			sub.RunUnique(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
