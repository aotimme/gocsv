package cmd

import (
	"fmt"
	"testing"
)

func TestRunTranspose(t *testing.T) {
	testCases := []struct {
		rows [][]string
	}{
		{[][]string{
			{"Name", "DataFox Intelligence, Inc."},
			{"Website", "www.datafox.com"},
		}},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ic, err := NewInputCsv("../test-files/simple.csv")
			if err != nil {
				t.Error("Unexpected error", err)
			}
			toc := new(testOutputCsv)
			sub := new(TransposeSubcommand)
			sub.RunTranspose(ic, toc)
			err = assertRowsEqual(tt.rows, toc.rows)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
