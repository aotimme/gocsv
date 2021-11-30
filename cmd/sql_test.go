package cmd

import (
	"fmt"
	"testing"
)

func TestRunSql(t *testing.T) {
	testCases := []struct {
		queryString string
		rows        [][]string
	}{
		{"SELECT * FROM [simple-sort] WHERE [Number] > 0", [][]string{
			[]string{"Number", "String"},
			[]string{"1", "One"},
			[]string{"2", "Two"},
			[]string{"2", "Another Two"},
		}},
		{"SELECT SUM([Number]) AS Total FROM [simple-sort]", [][]string{
			[]string{"Total"},
			[]string{"4"},
		}},
		{"SELECT [Number], COUNT(*) AS Count FROM [simple-sort] GROUP BY [Number] ORDER BY [Number] ASC", [][]string{
			[]string{"Number", "Count"},
			[]string{"-1", "1"},
			[]string{"1", "1"},
			[]string{"2", "2"},
		}},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ic, err := NewInputCsv("../test-files/simple-sort.csv")
			if err != nil {
				t.Error("Unexpected error", err)
			}
			toc := new(testOutputCsv)
			sub := new(SqlSubcommand)
			sub.queryString = tt.queryString
			sub.RunSql([]*InputCsv{ic}, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEscapeSqlName(t *testing.T) {
	testCases := []struct {
		inputName  string
		outputName string
	}{
		{"basic", "\"basic\""},
		{"single space", "\"single space\""},
		{"single'quote", "\"single'quote\""},
		{"square[]brackets", "\"square[]brackets\""},
		{"\"alreadyquoted\"", "\"\"\"alreadyquoted\"\"\""},
		{"middle\"quote", "\"middle\"\"quote\""},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			output := escapeSqlName(tt.inputName)
			if output != tt.outputName {
				t.Errorf("Expected %s but got %s", tt.outputName, output)
			}
		})
	}
}
