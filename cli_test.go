package habit_test

import (
	"bytes"
	"habit"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestNoArgsShowsAllHabits(t *testing.T) {
	t.Parallel()
	var args []string
	buffer := bytes.Buffer{}
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())
	habit.RunCLI(tmpFile.Name(), args, &buffer)

	want := "Habits:"
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("No arguments should print usage Message got: %s", got)
	}
}
func TestMoreThanOneArgShowsUsageHelp(t *testing.T) {
	t.Parallel()
	args := []string{"blah", "blah"}
	buffer := bytes.Buffer{}
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())

	want := "Usage"
	habit.RunCLI(tmpFile.Name(), args, &buffer)
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("No arguments should print usage Message got: %s", got)
	}
}

func TestOptionsButNoArgsShowsUsageHelp(t *testing.T) {
	t.Parallel()
	args := []string{"-f", "daily"}
	buffer := bytes.Buffer{}
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())

	want := "Usage"
	habit.RunCLI(tmpFile.Name(), args, &buffer)
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("No arguments should print usage Message got: %s", got)
	}
}

func TestOptionServerStartsAHTTPSERVER(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "localhost:8080"}
	buffer := bytes.Buffer{}
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())

	//TODO not the correct way to run a server... buffer is not needed?
	go habit.RunCLI(tmpFile.Name(), args, &buffer)
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Want Status 200, got: %d", resp.StatusCode)
	}
}
