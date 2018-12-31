package main

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
			[]string{"Number", "String"},
			[]string{"2", "Another Two"},
		}},
		{"0", [][]string{
			[]string{"Number", "String"},
		}},
		{"+1", [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Two"},
			[]string{"-1", "Minus One"},
			[]string{"2", "Another Two"},
		}},
		{"1000", [][]string{
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
