package main

import (
	"fmt"
	"testing"
)

func TestRunReplace(t *testing.T) {
	testCases := []struct {
		columnsString   string
		regex           string
		repl            string
		caseInsensitive bool
		rows            [][]string
	}{
		{"String", "Two", "Dos", false, [][]string{
			[]string{"Number", "String"},
			[]string{"1", "One"},
			[]string{"2", "Dos"},
			[]string{"-1", "Minus One"},
			[]string{"2", "Another Dos"},
		}},
		{"String", "^one", "UNO", true, [][]string{
			[]string{"Number", "String"},
			[]string{"1", "UNO"},
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
			sub := new(ReplaceSubcommand)
			sub.columnsString = tt.columnsString
			sub.regex = tt.regex
			sub.repl = tt.repl
			sub.caseInsensitive = tt.caseInsensitive
			sub.RunReplace(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
