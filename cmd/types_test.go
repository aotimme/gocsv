package cmd

import (
	"os"
	"testing"
	"time"
)

func TestTimeLayoutEnvVar(t *testing.T) {
	const layout = "06/01/02 3:04pm"
	os.Setenv("GOCSV_TIMELAYOUT", layout)

	s := "00/01/01 10:34pm"
	want := time.Date(2000, 1, 1, 22, 34, 0, 0, time.UTC)

	useTimeLayoutEnvVar()

	if got, err := ParseDatetime(s); err != nil || got != want {
		t.Errorf("ParseDatetime(%q) = %v, %v; want %v, nil", s, got, err, want)
	}

	if got, err := ParseDate(s); err != nil || got != want {
		t.Errorf("ParseDate(%q) = %v, %v; want %v, nil", s, got, err, want)
	}
}
