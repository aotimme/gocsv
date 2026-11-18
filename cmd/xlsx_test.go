package cmd

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestXlsxPadsTrailingBlankCells(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	for cell, value := range map[string]string{
		"A1": "name", "B1": "email", "C1": "notes",
		"A2": "Ann", "B2": "ann@example.com",
		"A3": "Ben", "B3": "ben@example.com", "C3": "VIP",
	} {
		if err := f.SetCellValue("Sheet1", cell, value); err != nil {
			t.Fatal(err)
		}
	}

	toc := new(testOutputCsv)
	writeRowsToOutputCsv(toc, f, "Sheet1")
	err := assertRowsEqual([][]string{
		{"name", "email", "notes"},
		{"Ann", "ann@example.com", ""},
		{"Ben", "ben@example.com", "VIP"},
	}, toc.rows)
	if err != nil {
		t.Error(err)
	}
}
