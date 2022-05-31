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
	store := habit.OpenMemoryStore()
	habit.RunCLI([]string{}, &buffer, &store)

	got := buffer.String()
	if !strings.Contains(got, "Usage") {
		t.Errorf("no arguments and no previous habits should print usage Message got:\n    %s", got)
	}
}

func TestNoArgsShowsAllHabitsWithExistingHabits(t *testing.T) {
	t.Parallel()
	buffer := bytes.Buffer{}
	store := habit.OpenMemoryStore()
	store.Habits = map[string]*habit.Habit{
		"piano": &habit.Habit{Name: "piano"},
	}
	habit.RunCLI([]string{}, &buffer, &store)

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
	habit.RunCLI(args, &buffer, nil)

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
	habit.RunCLI(args, &buffer, nil)
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
		store := habit.OpenMemoryStore()
		habit.RunCLI(tc.args, &buffer, &store)
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
	err := habit.RunCLI(args, &buffer, nil)

	if err == nil {
		t.Errorf("invalid frequency should return error. Got nil")
	}
}
