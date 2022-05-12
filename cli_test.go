package habit_test

import (
	"bytes"
	"habit"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	usage  = "Usage"
	habits = "Habits:"
)

func TestNoArgsShowsUsageHelpOnNoHabits(t *testing.T) {
	t.Parallel()
	var args []string
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())
	habit.RunCLI(tmpFile.Name(), args, &buffer)

	want := usage
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("No arguments and no previous habits should print usage Message got\n:  %s", got)
	}
}

func TestNoArgsShowsAllHabitsWithExistingHabits(t *testing.T) {
	t.Parallel()
	var args []string
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())
	writeTracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Interval: habit.WeeklyInterval,
			Streak:   1,
			DueDate:  time.Now().Add(habit.WeeklyInterval),
		},
	}

	writeFileStore := habit.NewFileStore(tmpFile.Name())
	writeFileStore.Write(&writeTracker)
	habit.RunCLI(tmpFile.Name(), args, &buffer)

	want := habits
	got := buffer.String()
	if !strings.Contains(got, want) {
		// if habits exist in store, a summary of all habits should be displayed
		t.Errorf("No arguments and previous habits should print a summary of all habits got\n:  %s", got)
	}
}
func TestMoreThanOneArgShowsUsageHelp(t *testing.T) {
	t.Parallel()
	args := []string{"blah", "blah"}
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())

	want := usage
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
	tmpFile := CreateTmpFile()
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
	args := []string{"-s", address}
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())

	go habit.RunCLI(tmpFile.Name(), args, &buffer)
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Want Status 200, got: %d", resp.StatusCode)
	}
}
