package habit_test

import (
	"bytes"
	"habit"
	"strings"
	"testing"
)

func TestNoArgsShowsAllHabits(t *testing.T) {
	var args []string
	buffer := bytes.Buffer{}
	want := "Habits:"
	habit.RunCLI(args, &buffer)
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("No arguments should print usage message got: %s", got)
	}
}
func TestMoreThanOneArgShowsUsageHelp(t *testing.T) {
	args := []string{"blah", "blah"}
	buffer := bytes.Buffer{}
	want := "Usage"
	habit.RunCLI(args, &buffer)
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("No arguments should print usage message got: %s", got)
	}
}

func TestOptionsButNoArgsShowsUsageHelp(t *testing.T) {
	args := []string{"-f", "daily"}
	buffer := bytes.Buffer{}
	want := "Usage"
	habit.RunCLI(args, &buffer)
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("No arguments should print usage message got: %s", got)
	}
}
