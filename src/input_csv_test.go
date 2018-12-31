package main

import (
	"testing"
)

func TestNewInputCsv(t *testing.T) {
	ic, err := NewInputCsv("../test-files/simple.csv")
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if ic.filename != "../test-files/simple.csv" {
		t.Error("Unexpected filename", ic.filename)
	}
	if ic.hasBom {
		t.Error("Expected no BOM")
	}
	err = ic.Close()
	if err != nil {
		t.Error("Unexpected error from close", err)
	}
}

func TestReadAll(t *testing.T) {
	ic, err := NewInputCsv("../test-files/simple.csv")
	if err != nil {
		t.Error("Unexpected error", err)
	}
	rows, err := ic.ReadAll()
	if err != nil {
		t.Error("Unexpected error reading all", err)
	}
	if len(rows) != 2 {
		t.Error("Expected 2 rows but got", len(rows))
	}
	expected := [][]string{
		[]string{"Name", "Website"},
		[]string{"DataFox Intelligence, Inc.", "www.datafox.com"},
	}
	for i, row := range expected {
		for j, cell := range row {
			if cell != rows[i][j] {
				t.Error("Expected", cell, "at", i, j, "but got", rows[i][j])
			}
		}
	}
}

func TestGetInputCsvs(t *testing.T) {
	testCases := []struct {
		description  string
		filenames    []string
		numInputCsvs int
		numCsvs      int
	}{
		{"1 input without filename", []string{}, 1, 1},
		{"1 input with stdin", []string{"-"}, 1, 1},
		{"1 input with filename", []string{"../test-files/simple.csv"}, 1, 1},
		{"2 inputs with one filename", []string{"../test-files/simple.csv"}, 2, 2},
		{"2 inputs with filename", []string{"../test-files/simple.csv", "../test-files/simple-bom.csv"}, 2, 2},
		{"many inputs with no filename", []string{}, -1, 1},
		{"many inputs with 1 filename", []string{"../test-files/simple.csv"}, -1, 1},
		{"many inputs with 2 filenames", []string{"../test-files/simple.csv", "../test-files/simple-bom.csv"}, -1, 2},
		{"many inputs with stdin and 2 filenames", []string{"-", "../test-files/simple.csv", "../test-files/simple-bom.csv"}, -1, 3},
	}
	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			inputCsvs, err := GetInputCsvs(tt.filenames, tt.numInputCsvs)
			if err != nil {
				t.Error("Unexpected error", err)
			}
			if len(inputCsvs) != tt.numCsvs {
				t.Errorf("Expected %q CSVs but got %q", tt.numCsvs, len(inputCsvs))
			}
		})
	}
}

func TestGetInputCsvsErrors(t *testing.T) {
	testCases := []struct {
		description  string
		filenames    []string
		numInputCsvs int
	}{
		{"1 input with stdin and filename", []string{"-", "../test-files/simple.csv"}, 1},
		{"1 input with multiple filenames", []string{"../test-files/simple.csv", "../test-files/simple-bom.csv"}, 1},
		{"2 inputs with no filenames", []string{}, 2},
		{"2 inputs with stdin filename", []string{"-"}, 2},
		{"2 inputs with 3 filenames", []string{"-", "../test-files/simple.csv", "../test-files/simple-bom.csv"}, 2},
	}
	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			_, err := GetInputCsvs(tt.filenames, tt.numInputCsvs)
			if err == nil {
				t.Error("Expected error but got nil")
			}
		})
	}
}
