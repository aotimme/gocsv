package cmd

import (
	"fmt"
	"testing"
)

func TestRunBehead(t *testing.T) {
	testCases := []struct {
		numHeaders int
		rows       [][]string
	}{
		{1, [][]string{
			{"1", "One"},
			{"2", "Two"},
			{"-1", "Minus One"},
			{"2", "Another Two"},
		}},
		{2, [][]string{
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
			sub := new(BeheadSubcommand)
			sub.numHeaders = tt.numHeaders
			sub.RunBehead(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
