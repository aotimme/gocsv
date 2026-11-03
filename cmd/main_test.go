package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Run the real command boundary in a subprocess so tests observe os.Exit and
// distinguish normal error reporting from debug panics.
func TestMain(m *testing.M) {
	if os.Getenv("GOCSV_TEST_COMMAND") != "1" {
		os.Exit(m.Run())
	}
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = append([]string{"gocsv"}, os.Args[i+1:]...)
			break
		}
	}
	if os.Getenv("GOCSV_TEST_CLOSE_STDOUT") == "1" {
		os.Stdout.Close()
	}
	Main()
	os.Exit(0)
}

func commandForTest(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command(executable, append([]string{"--"}, args...)...)
	command.Env = append(os.Environ(), "GOCSV_TEST_COMMAND=1", "GOCSV_FLUSH_EVERY_ROW=", "GOCSV_DELIMITER=")
	return command
}

func TestCommandErrorsFlushBufferedRows(t *testing.T) {
	for _, numRows := range []int{1, 2000} {
		for _, debug := range []bool{false, true} {
			t.Run(strconv.Itoa(numRows)+"/debug="+strconv.FormatBool(debug), func(t *testing.T) {
				var input, want strings.Builder
				input.WriteString("Name,Value\n")
				want.WriteString("Name\n")
				for i := 0; i < numRows; i++ {
					name := strconv.Itoa(i)
					input.WriteString(name + ",1\n")
					want.WriteString(name + "\n")
				}
				input.WriteString("bad\n")
				command := commandForTest(t, "select", "-c", "1", "--debug="+strconv.FormatBool(debug))
				command.Stdin = strings.NewReader(input.String())
				var stdout, stderr bytes.Buffer
				command.Stdout, command.Stderr = &stdout, &stderr
				err := command.Run()
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) || (!debug && exitErr.ExitCode() != 1) {
					t.Fatalf("expected command failure, got %v; stderr: %s", err, &stderr)
				}
				if stdout.String() != want.String() {
					t.Fatalf("buffered rows lost: got %d bytes, want %d", stdout.Len(), want.Len())
				}
				if !strings.Contains(stderr.String(), "wrong number of fields") {
					t.Fatalf("missing parse error: %s", &stderr)
				}
				if strings.Contains(stderr.String(), "panic:") != debug {
					t.Fatalf("unexpected error reporting for debug=%v: %s", debug, &stderr)
				}
			})
		}
	}
}

func TestCommandReportsFinalFlushError(t *testing.T) {
	command := commandForTest(t, "select", "-c", "1")
	command.Env = append(command.Env, "GOCSV_TEST_CLOSE_STDOUT=1")
	command.Stdin = strings.NewReader("Name,Value\na,1\n")
	output, err := command.CombinedOutput()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("expected exit status 1, got %v: %s", err, output)
	}
	if !strings.HasPrefix(string(output), "Error: ") || !strings.Contains(string(output), "file already closed") {
		t.Fatalf("unexpected flush error: %s", output)
	}
}

func TestSplitErrorFlushesCurrentFile(t *testing.T) {
	base := filepath.Join(t.TempDir(), "out")
	command := commandForTest(t, "split", "--max-rows", "1", "--filename-base", base)
	command.Stdin = strings.NewReader("Name,Value\na,1\nb,2\nbad\n")
	output, err := command.CombinedOutput()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("expected exit status 1, got %v: %s", err, output)
	}
	for i, want := range []string{"Name,Value\na,1\n", "Name,Value\nb,2\n"} {
		got, err := os.ReadFile(base + "-" + strconv.Itoa(i+1) + ".csv")
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("part %d: got %q, want %q", i+1, got, want)
		}
	}
}

func TestSplitClosesCompletedFiles(t *testing.T) {
	shell, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("requires a shell with ulimit")
	}
	base := filepath.Join(t.TempDir(), "out")
	command := commandForTest(t, "split", "--max-rows", "1", "--filename-base", base)
	// More output parts than available descriptors catches retained open files.
	command.Args = append([]string{shell, "-c", `ulimit -n 64 || exit 99; exec "$@"`, "sh"}, command.Args...)
	command.Path = shell
	command.Stdin = strings.NewReader("Name\n" + strings.Repeat("a\n", 128))
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("split failed with a low descriptor limit: %v: %s", err, output)
	}
	for i := 1; i <= 128; i++ {
		got, err := os.ReadFile(base + "-" + strconv.Itoa(i) + ".csv")
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != "Name\na\n" {
			t.Fatalf("part %d: unexpected output %q", i, got)
		}
	}
}
