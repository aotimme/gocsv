package main

import (
	"fmt"
	"testing"
)

func TestSortCsv(t *testing.T) {
	testCases := []struct {
		columns     string
		reverse     bool
		noInference bool
		rows        [][]string
	}{
		{"Number", false, false, [][]string{
			[]string{"Number", "String"},
			[]string{"-1", "Minus One"},
			[]string{"1", "One"},
			[]string{"2", "Another Two"},
			[]string{"2", "Two"},
		}},
		{"Number", true, false, [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Another Two"},
			[]string{"2", "Two"},
			[]string{"1", "One"},
			[]string{"-1", "Minus One"},
		}},
		{"Number", false, true, [][]string{
			[]string{"Number", "String"},
			[]string{"-1", "Minus One"},
			[]string{"1", "One"},
			[]string{"2", "Another Two"},
			[]string{"2", "Two"},
		}},
		{"Number", true, true, [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Another Two"},
			[]string{"2", "Two"},
			[]string{"1", "One"},
			[]string{"-1", "Minus One"},
		}},
		{"String", false, false, [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Another Two"},
			[]string{"-1", "Minus One"},
			[]string{"1", "One"},
			[]string{"2", "Two"},
		}},
		{"Number,String", true, true, [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Two"},
			[]string{"2", "Another Two"},
			[]string{"1", "One"},
			[]string{"-1", "Minus One"},
		}},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ic, err := NewInputCsv("../test-files/simple-sort.csv")
			if err != nil {
				t.Error("Unexpected error", err)
			}
			toc := new(testOutputCsv)
			sub := new(SortSubcommand)
			sub.columnsString = tt.columns
			sub.reverse = tt.reverse
			sub.noInference = tt.noInference
			sub.SortCsv(ic, toc)
			if len(tt.rows) != len(toc.rows) {
				t.Errorf("Expected %d rows but got %d", len(tt.rows), len(toc.rows))
			}
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
