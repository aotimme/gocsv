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
			[]string{"Number", "String", ""},
			[]string{"1", "One", "Row 1: 1 (One)"},
			[]string{"2", "Two", "Row 2: 2 (Two)"},
			[]string{"-1", "Minus One", "Row 3: -1 (Minus One)"},
			[]string{"2", "Another Two", "Row 4: 2 (Another Two)"},
		}},
		{"Long Version", "Row {{.index}}: {{.Number}} ({{.String}})", false, [][]string{
			[]string{"Number", "String", "Long Version"},
			[]string{"1", "One", "Row 1: 1 (One)"},
			[]string{"2", "Two", "Row 2: 2 (Two)"},
			[]string{"-1", "Minus One", "Row 3: -1 (Minus One)"},
			[]string{"2", "Another Two", "Row 4: 2 (Another Two)"},
		}},
		{"", "Row {{.index}}: {{.Number}} ({{.String}})", true, [][]string{
			[]string{"", "Number", "String"},
			[]string{"Row 1: 1 (One)", "1", "One"},
			[]string{"Row 2: 2 (Two)", "2", "Two"},
			[]string{"Row 3: -1 (Minus One)", "-1", "Minus One"},
			[]string{"Row 4: 2 (Another Two)", "2", "Another Two"},
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
