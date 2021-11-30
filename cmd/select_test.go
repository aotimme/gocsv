package cmd

import (
	"fmt"
	"testing"
)

func TestRunSelect(t *testing.T) {
	testCases := []struct {
		columnsString string
		exclude       bool
		rows          [][]string
	}{
		{"String", false, [][]string{
			{"String"},
			{"One"},
			{"Two"},
			{"Minus One"},
			{"Another Two"},
		}},
		{"Number", true, [][]string{
			{"String"},
			{"One"},
			{"Two"},
			{"Minus One"},
			{"Another Two"},
		}},
		{"String,Number", false, [][]string{
			{"String", "Number"},
			{"One", "1"},
			{"Two", "2"},
			{"Minus One", "-1"},
			{"Another Two", "2"},
		}},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ic, err := NewInputCsv("../test-files/simple-sort.csv")
			if err != nil {
				t.Error("Unexpected error", err)
			}
			toc := new(testOutputCsv)
			sub := new(SelectSubcommand)
			sub.columnsString = tt.columnsString
			sub.exclude = tt.exclude
			sub.RunSelect(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
