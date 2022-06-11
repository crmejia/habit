package habit_test

import (
	"bytes"
	"github.com/crmejia/habit"
	"strings"
	"testing"
)

func TestRunCLIShowsUsageHelpNoArgs(t *testing.T) {
	t.Parallel()
	buffer := bytes.Buffer{}
	habit.RunCLI([]string{}, &buffer)

	got := buffer.String()
	if !strings.Contains(got, "Usage") {
		t.Errorf("no arguments and no previous habits should print usage Message got:\n    %s", got)
	}
}

func TestRunCLIAllShowsCorrectly(t *testing.T) {
	t.Parallel()
	args := []string{"all"}
	buffer := bytes.Buffer{}
	habit.RunCLI(args, &buffer)

	want := "Habits:"
	got := buffer.String()
	if !strings.Contains(got, want) {
		// if habits exist in store, a summary of all habits should be displayed
		t.Errorf("habit all should print a summary of all habits got\n:  %s", got)
	}
}

func TestRunCLIShowsUsageHelpMoreThanOneArg(t *testing.T) {
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

func TestRunCLIShowsUsageHelpOptionsButNoArgs(t *testing.T) {
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

func TestRunCLIShowsErrorUsageHelpWrongOptions(t *testing.T) {
	t.Parallel()
	args := []string{"-g", "gibberish"}
	buffer := bytes.Buffer{}

	habit.RunCLI(args, &buffer)
	got := buffer.String()
	if !strings.Contains(got, "Usage") {
		t.Errorf("only options and no arguments should print usage Message got: %s", got)
	}
}

func TestRunClIShowNewHabitMessageNewHabit(t *testing.T) {
	t.Parallel()
	buffer := bytes.Buffer{}
	tmpDir := t.TempDir()
	args := []string{"-d", tmpDir, "piano"}

	habit.RunCLI(args, &buffer)
	want := "Good luck with your new habit"
	got := buffer.String()
	if !strings.Contains(got, want) {
		t.Errorf("new habit should print streak message. Got:\n  %s", got)
	}
}

func TestRunCLIShowsErrorUsageHelpInvalidFrequency(t *testing.T) {
	t.Parallel()
	args := []string{"-f", "yellow", "piano"}
	buffer := bytes.Buffer{}
	habit.RunCLI(args, &buffer)

	got := buffer.String()

	if !strings.Contains(got, "unknown frequency:") {
		t.Errorf("Invalid frecuency should print error message got: %s", got)
	}
	if !strings.Contains(got, "Usage") {
		t.Errorf("Invalid frecuency should print usage Message got: %s", got)
	}
}

func TestRunCLIShowsErrorUsageHelpEmptyFrequency(t *testing.T) {
	t.Parallel()
	args := []string{"-f", "", "piano"}
	buffer := bytes.Buffer{}
	habit.RunCLI(args, &buffer)

	got := buffer.String()

	if !strings.Contains(got, "habit frequency cannot be empty") {
		t.Errorf("Empty frecuency should print error message got: %s", got)
	}
	if !strings.Contains(got, "Usage") {
		t.Errorf("Empty frecuency should print usage Message got: %s", got)
	}
}

func TestRunCLIShowsErrorUsageHelpInvalidStoreType(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "cloud", "piano"}
	buffer := bytes.Buffer{}
	habit.RunCLI(args, &buffer)

	got := buffer.String()

	if !strings.Contains(got, "unknown store type") {
		t.Errorf("Invalid frecuency should print error message got: %s", got)
	}
	if !strings.Contains(got, "Usage") {
		t.Errorf("Invalid frecuency should print usage Message got: %s", got)
	}
}
