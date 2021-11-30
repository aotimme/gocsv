package cmd

import (
	"fmt"
	"testing"
)

func TestRunTail(t *testing.T) {
	testCases := []struct {
		numRowsStr string
		rows       [][]string
	}{
		{"1", [][]string{
			{"Number", "String"},
			{"2", "Another Two"},
		}},
		{"0", [][]string{
			{"Number", "String"},
		}},
		{"+1", [][]string{
			{"Number", "String"},
			{"2", "Two"},
			{"-1", "Minus One"},
			{"2", "Another Two"},
		}},
		{"1000", [][]string{
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
			sub := new(TailSubcommand)
			sub.numRowsStr = tt.numRowsStr
			sub.RunTail(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
