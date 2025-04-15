package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestRunSql(t *testing.T) {
	testCases := []struct {
		queryString string
		rows        [][]string
	}{
		{"SELECT * FROM [simple-sort] WHERE [Number] > 0", [][]string{
			{"Number", "String"},
			{"1", "One"},
			{"2", "Two"},
			{"2", "Another Two"},
		}},
		{"SELECT SUM([Number]) AS Total FROM [simple-sort]", [][]string{
			{"Total"},
			{"4"},
		}},
		{"SELECT [Number], COUNT(*) AS Count FROM [simple-sort] GROUP BY [Number] ORDER BY [Number] ASC", [][]string{
			{"Number", "Count"},
			{"-1", "1"},
			{"1", "1"},
			{"2", "2"},
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

func TestTxtar(t *testing.T) {
	const (
		inPrefix   = "in: "
		testPrefix = "test: "
		wantPrefix = "want: "
	)
	var (
		join     func(name string) (tempPath string)
		chTmpDir func()
	)
	{
		tmpPath := t.TempDir()
		join = func(name string) (tempPath string) {
			return filepath.Join(tmpPath, name)
		}
		chTmpDir = func() { os.Chdir(tmpPath) }

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("could not get working dir: %v", err)
		}
		t.Cleanup(func() {
			err := os.Chdir(cwd)
			if err != nil {
				t.Fatalf("could not change back to orignal working dir %s: %v", cwd, err)
			}
		})
	}

	a, err := txtar.ParseFile("./testdata/sql.txt")
	if err != nil {
		t.Fatal(err)
	}

	testPairs := []txtar.File{}
	for _, aFile := range a.Files {
		name := aFile.Name
		if !strings.HasPrefix(name, inPrefix) {
			testPairs = append(testPairs, aFile)
			continue
		}
		name = strings.TrimPrefix(name, inPrefix)
		path := join(name)
		if err := os.WriteFile(path, aFile.Data, 0644); err != nil {
			t.Fatal(err)
		}
		if strings.HasSuffix(name, ".csv") {
			if err := preProcessCSV(path); err != nil {
				t.Fatalf("got non-nil error for preProcessCSV(%s); %v", name, err)
			}
		}
	}

	chTmpDir()

	for i := 0; i < len(testPairs); i += 2 {
		testFile := testPairs[i]
		wantFile := testPairs[i+1]
		if !strings.HasPrefix(testFile.Name, testPrefix) ||
			!strings.HasPrefix(wantFile.Name, wantPrefix) {
			t.Fatalf("got test-file=%q want-file=%q; want \"test: ...\", \"want: ...\"",
				testFile.Name, wantFile.Name)
		}
		testName := strings.TrimPrefix(testFile.Name, testPrefix)
		wantName := strings.TrimPrefix(wantFile.Name, wantPrefix)
		if testName != wantName {
			t.Fatalf("got test-name=%s want-name=%s; want test-name==want-name", testName, wantName)
		}

		t.Run(testName, func(t *testing.T) {
			var (
				args   []string
				subcmd SqlSubcommand
			)

			err := json.Unmarshal(testFile.Data, &args)
			if err != nil {
				t.Fatalf("expected test data to be an array of string values: %v", err)
			}

			fs := flag.NewFlagSet(subcmd.Name(), flag.ExitOnError)
			subcmd.SetFlags(fs)
			err = fs.Parse(args)
			if err != nil {
				t.Fatalf("could not parse args: %v", err)
			}

			buf := bytes.Buffer{}

			inputCsvs := GetInputCsvsOrPanic(fs.Args(), -1)
			outputCsv := &OutputCsv{
				csvWriter: csv.NewWriter(&buf),
			}
			subcmd.RunSql(inputCsvs, outputCsv)

			b, err := postProcessCSV(buf.Bytes())
			if err != nil {
				t.Fatalf("could post-process CSV: %v", err)
			}
			got := strings.TrimSpace(string(b))
			want := strings.TrimSpace(string(wantFile.Data))
			if got != want {
				t.Errorf("\n got: \n%q\nwant: \n%q", got, want)
			}
		})
	}
}

// preProcessCSV modifies and overwrites path
// to de-prettify the CSV.
func preProcessCSV(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	f, err = os.Create(path)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}
	writer.Flush()
	if err = writer.Error(); err != nil {
		return err
	}

	return nil
}

// postProcessCSV prettifies data.
func postProcessCSV(data []byte) ([]byte, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var widths []int
	for _, record := range records[0] {
		widths = append(widths, len(record))
	}
	for _, record := range records[1:] {
		for i, field := range record {
			widths[i] = max(len(field), widths[i])
		}
	}

	buf := &bytes.Buffer{}
	for _, record := range records {
		line := record[0]
		sep := ", "
		for i := 1; i < len(record); i++ {
			field := record[i]
			field = strings.Repeat(" ", widths[i-1]-len(record[i-1])) + field
			line += sep + field
		}
		if _, err := buf.WriteString(line + "\n"); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
