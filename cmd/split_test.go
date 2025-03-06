package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestSplit(t *testing.T) {
	tempdir := t.TempDir()

	archive, err := txtar.ParseFile(filepath.Join("testdata", "split.txt"))
	if err != nil {
		t.Fatal(err)
	}
	files := archive.Files

	// for split, assert single input file at top of list,
	// and create it for the split subcommand to read
	input := files[0]
	if input.Name != "input.csv" {
		t.Fatalf("got first file %s; want input.csv", input.Name)
	}
	// need to change dir to test that calling split without
	// a basename works
	os.Chdir(tempdir)
	f, err := os.Create(input.Name)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write(input.Data)
	if err != nil {
		t.Fatal(err)
	}

	// iterate remaining files, breaking them up into groups of
	// a single test file and some number of subsequent want files
	const (
		testPrefix = "test: "
		wantPrefix = "want: "
	)

	files = files[1:]
	for i := 0; i < len(files); {
		if !strings.HasPrefix(files[i].Name, testPrefix) {
			t.Fatalf("got first file %s in group; want test: ...", files[0].Name)
		}
		test := files[i]
		i++

		wantFiles := make([]txtar.File, 0)
		for i < len(files) &&
			strings.HasPrefix(files[i].Name, wantPrefix) {
			wantFiles = append(wantFiles, files[i])
			i++
		}
		if len(wantFiles) == 0 {
			t.Fatal("got 0 want-files; want some number of want-files")
		}

		t.Run(strings.TrimPrefix(test.Name, testPrefix), func(t *testing.T) {
			args := strings.Fields(trimb(test.Data))
			args = append(args, input.Name)

			sc := SplitSubcommand{}
			fs := flag.NewFlagSet(sc.Name(), flag.ExitOnError)
			sc.SetFlags(fs)
			err := fs.Parse(args)
			if err != nil {
				t.Fatalf("could not parse args %q: %v", args, err)
			}
			sc.Run(fs.Args())

			for _, file := range wantFiles {
				name := strings.TrimPrefix(file.Name, wantPrefix)
				want := file.Data

				got, err := os.ReadFile(name)
				if err != nil {
					t.Fatal(err)
				}

				if trimb(got) != trimb(want) {
					t.Errorf("file %s\n got: %q\nwant: %q", name, trimb(got), trimb(want))
				}
			}
		})
	}
}

// trimb converts b to a string and calls [strings.TrimSpace]
// on it.
func trimb(b []byte) string {
	return strings.TrimSpace(string(b))
}
