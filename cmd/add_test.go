package cmd

import (
	"fmt"
	"testing"
)

func TestRunAdd(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		prepend  bool
		rows     [][]string
	}{
		{"", "Row {{.index}}: {{.Number}} ({{.String}})", false, [][]string{
			{"Number", "String", ""},
			{"1", "One", "Row 1: 1 (One)"},
			{"2", "Two", "Row 2: 2 (Two)"},
			{"-1", "Minus One", "Row 3: -1 (Minus One)"},
			{"2", "Another Two", "Row 4: 2 (Another Two)"},
		}},
		{"Long Version", "Row {{.index}}: {{.Number}} ({{.String}})", false, [][]string{
			{"Number", "String", "Long Version"},
			{"1", "One", "Row 1: 1 (One)"},
			{"2", "Two", "Row 2: 2 (Two)"},
			{"-1", "Minus One", "Row 3: -1 (Minus One)"},
			{"2", "Another Two", "Row 4: 2 (Another Two)"},
		}},
		{"", "Row {{.index}}: {{.Number}} ({{.String}})", true, [][]string{
			{"", "Number", "String"},
			{"Row 1: 1 (One)", "1", "One"},
			{"Row 2: 2 (Two)", "2", "Two"},
			{"Row 3: -1 (Minus One)", "-1", "Minus One"},
			{"Row 4: 2 (Another Two)", "2", "Another Two"},
		}},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ic, err := NewInputCsv("../test-files/simple-sort.csv")
			if err != nil {
				t.Error("Unexpected error", err)
			}
			toc := new(testOutputCsv)
			sub := new(AddSubcommand)
			sub.name = tt.name
			sub.template = tt.template
			sub.prepend = tt.prepend
			sub.RunAdd(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
