package main

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
			[]string{"String"},
			[]string{"One"},
			[]string{"Two"},
			[]string{"Minus One"},
			[]string{"Another Two"},
		}},
		{"Number", true, [][]string{
			[]string{"String"},
			[]string{"One"},
			[]string{"Two"},
			[]string{"Minus One"},
			[]string{"Another Two"},
		}},
		{"String,Number", false, [][]string{
			[]string{"String", "Number"},
			[]string{"One", "1"},
			[]string{"Two", "2"},
			[]string{"Minus One", "-1"},
			[]string{"Another Two", "2"},
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
