package cmd

import (
	"fmt"
	"testing"
)

func TestRunFilter(t *testing.T) {
	testCases := []struct {
		columnsString   string
		exclude         bool
		regex           string
		equals          string
		caseInsensitive bool
		gtStr           string
		gteStr          string
		ltStr           string
		lteStr          string
		rows            [][]string
	}{
		// gt
		{"Number", false, "", "", false, "1", "", "", "", [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Two"},
			[]string{"2", "Another Two"},
		}},
		// gte
		{"Number", false, "", "", false, "", "1", "", "", [][]string{
			[]string{"Number", "String"},
			[]string{"1", "One"},
			[]string{"2", "Two"},
			[]string{"2", "Another Two"},
		}},
		// lt
		{"Number", false, "", "", false, "", "", "1", "", [][]string{
			[]string{"Number", "String"},
			[]string{"-1", "Minus One"},
		}},
		// lte
		{"Number", false, "", "", false, "", "", "", "1", [][]string{
			[]string{"Number", "String"},
			[]string{"1", "One"},
			[]string{"-1", "Minus One"},
		}},
		// equals
		{"String", false, "", "Two", false, "", "", "", "", [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Two"},
		}},
		// regex
		{"String", false, "[tT]wo", "", false, "", "", "", "", [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Two"},
			[]string{"2", "Another Two"},
		}},
		// regex exclude
		{"String", true, "[oO]ne", "", false, "", "", "", "", [][]string{
			[]string{"Number", "String"},
			[]string{"2", "Two"},
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
			sub := new(FilterSubcommand)
			sub.columnsString = tt.columnsString
			sub.exclude = tt.exclude
			sub.regex = tt.regex
			sub.equals = tt.equals
			sub.caseInsensitive = tt.caseInsensitive
			sub.gtStr = tt.gtStr
			sub.gteStr = tt.gteStr
			sub.ltStr = tt.ltStr
			sub.lteStr = tt.lteStr
			sub.RunFilter(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
