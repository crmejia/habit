package habit_test

import (
	"bytes"
	"habit"
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
		t.Errorf("No arguments should print usage message got: %s", got)
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
		t.Errorf("No arguments should print usage message got: %s", got)
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
		t.Errorf("No arguments should print usage message got: %s", got)
	}
}
