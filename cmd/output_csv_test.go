package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestOutputCsvFlushEveryRow(t *testing.T) {
	testCases := []struct {
		value string
		flush bool
	}{
		{"unset", false}, {"", false},
		{"1", true}, {"t", true}, {"T", true},
		{"true", true}, {"TRUE", true}, {"True", true},
		{"0", false}, {"f", false}, {"F", false},
		{"false", false}, {"FALSE", false}, {"False", false},
	}
	for _, tt := range testCases {
		t.Run(tt.value, func(t *testing.T) {
			t.Setenv("GOCSV_DELIMITER", "")
			t.Setenv("GOCSV_FLUSH_EVERY_ROW", tt.value)
			if tt.value == "unset" {
				if err := os.Unsetenv("GOCSV_FLUSH_EVERY_ROW"); err != nil {
					t.Fatal(err)
				}
			}
			file, err := os.CreateTemp(t.TempDir(), "output-*.csv")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			output := NewOutputCsvFromFile(file)
			want := ""
			for _, row := range [][]string{{"Name", "Value"}, {"a,b", "1"}} {
				if err := output.Write(row); err != nil {
					t.Fatal(err)
				}
				got, err := os.ReadFile(file.Name())
				if err != nil {
					t.Fatal(err)
				}
				if tt.flush && len(got) <= len(want) {
					t.Fatalf("row was not flushed: %q", got)
				}
				if !tt.flush && len(got) != 0 {
					t.Fatalf("expected buffered output, got %q", got)
				}
				want = string(got)
			}
			output.FlushOrPanic()
			got, err := os.ReadFile(file.Name())
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != "Name,Value\n\"a,b\",1\n" {
				t.Fatalf("unexpected output after final flush: %q", got)
			}
		})
	}
}

func TestOutputCsvRejectsInvalidFlushEveryRow(t *testing.T) {
	debug := DEBUG
	DEBUG = true
	defer func() { DEBUG = debug }()

	for _, value := range []string{"invalid", "treu", "2", "yes", " true", " "} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("GOCSV_FLUSH_EVERY_ROW", value)
			defer func() {
				recovered := recover()
				err, ok := recovered.(error)
				if !ok {
					t.Fatalf("expected configuration error, got %v", recovered)
				}
				if !strings.Contains(err.Error(), "invalid GOCSV_FLUSH_EVERY_ROW value") {
					t.Fatalf("unexpected configuration error: %v", err)
				}
			}()
			NewOutputCsv()
		})
	}
}

func TestOutputCsvFlushOrPanicReportsWriteError(t *testing.T) {
	t.Setenv("GOCSV_FLUSH_EVERY_ROW", "")
	t.Setenv("GOCSV_DELIMITER", "")
	debug := DEBUG
	DEBUG = true
	defer func() { DEBUG = debug }()

	file, err := os.CreateTemp(t.TempDir(), "output-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	output := NewOutputCsvFromFile(file)
	if err := output.Write([]string{"buffered row"}); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		err, ok := recover().(error)
		if !ok || !errors.Is(err, os.ErrClosed) {
			t.Fatalf("expected closed-file error from deferred flush, got %v", err)
		}
	}()
	defer output.FlushOrPanic()
}

func TestSelectFlushesBufferedOutput(t *testing.T) {
	t.Setenv("GOCSV_FLUSH_EVERY_ROW", "")
	t.Setenv("GOCSV_DELIMITER", "")
	file, err := os.CreateTemp(t.TempDir(), "stdout-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	stdout := os.Stdout
	os.Stdout = file
	defer func() { os.Stdout = stdout }()

	sub := SelectSubcommand{columnsString: "1"}
	sub.Run([]string{"../test-files/simple.csv"})
	got, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "Name\n\"DataFox Intelligence, Inc.\"\n" {
		t.Fatalf("unexpected command output: %q", got)
	}
}

func TestXlsxFlushesBufferedOutput(t *testing.T) {
	t.Setenv("GOCSV_FLUSH_EVERY_ROW", "")
	t.Setenv("GOCSV_DELIMITER", "")
	workbook := excelize.NewFile()
	defer workbook.Close()
	if err := workbook.SetCellValue("Sheet1", "A1", "Name"); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	ConvertXlsxSheetToDirectory(workbook, dir, "Sheet1")
	got, err := os.ReadFile(filepath.Join(dir, "Sheet1.csv"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "Name\n" {
		t.Fatalf("unexpected sheet output: %q", got)
	}
}
