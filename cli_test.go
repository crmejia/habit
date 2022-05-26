package habit_test

import (
	"bytes"
	"habit"
	"strings"
	"testing"
)

const (
	usage  = "Usage"
	habits = "Habits:"
	streak = "streak"
)

func TestNoArgsShowsUsageHelpOnNoHabits(t *testing.T) {
	t.Parallel()
	buffer := bytes.Buffer{}
	habit.RunCLI([]string{}, &buffer)

	got := buffer.String()
	if !strings.Contains(got, "Usage") {
		t.Errorf("no arguments and no previous habits should print usage Message got:\n    %s", got)
	}
}

func TestNoArgsShowsAllHabitsWithExistingHabits(t *testing.T) {
	t.Parallel()
	t.Skip()
	buffer := bytes.Buffer{}

	habit.RunCLI([]string{}, &buffer)

	want := "Habits:"
	got := buffer.String()
	if !strings.Contains(got, want) {
		// if habits exist in store, a summary of all habits should be displayed
		t.Errorf("no arguments and previous habits should print a summary of all habits got\n:  %s", got)
	}
}

func TestMoreThanOneArgShowsUsageHelp(t *testing.T) {
	t.Parallel()
	args := []string{"blah", "blah"}
	buffer := bytes.Buffer{}
	habit.RunCLI(args, &buffer)

	want := "Usage"
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("too many arguments should print usage Message got: %s", got)
	}
}

func TestOptionsButNoArgsShowsUsageHelp(t *testing.T) {
	t.Parallel()
	args := []string{"-f", "daily"}
	buffer := bytes.Buffer{}

	want := "Usage"
	habit.RunCLI(args, &buffer)
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("only options and no arguments should print usage Message got: %s", got)
	}
}

func TestNewHabitShowNewHabitMessage(t *testing.T) {
	t.Parallel()
	buffer := bytes.Buffer{}
	testCases := []struct {
		args []string
		want string
	}{
		{args: []string{"piano"}, want: "again tomorrow."},
		{args: []string{"-f", "weekly", "piano"}, want: "again in a week."},
	}

	for _, tc := range testCases {
		habit.RunCLI(tc.args, &buffer)
		got := buffer.String()
		if !strings.Contains(got, tc.want) {
			t.Errorf("new habit should print streak message. Got:\n  %s", got)
		}
		buffer.Truncate(0)
	}
}

func TestNewHabitInvalidFrequencyReturnsError(t *testing.T) {
	t.Parallel()
	args := []string{"-f", "yellow", "piano"}
	buffer := bytes.Buffer{}
	err := habit.RunCLI(args, &buffer)

	if err == nil {
		t.Errorf("invalid frequency should return error. Got nil")
	}
}
