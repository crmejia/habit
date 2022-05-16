package habit_test

import (
	"bytes"
	"fmt"
	"github.com/phayes/freeport"
	"habit"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	usage  = "Usage"
	habits = "Habits:"
	streak = "streak"
)

func TestNoArgsShowsUsageHelpOnNoHabits(t *testing.T) {
	t.Parallel()
	var args []string
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile(t)

	habit.RunCLI(tmpFile.Name(), args, &buffer)

	want := usage
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("no arguments and no previous habits should print usage Message got\n:  %s", got)
	}
}

func TestNoArgsShowsAllHabitsWithExistingHabits(t *testing.T) {
	t.Parallel()
	var args []string
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile(t)

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
		t.Errorf("no arguments and previous habits should print a summary of all habits got\n:  %s", got)
	}
}

func TestMoreThanOneArgShowsUsageHelp(t *testing.T) {
	t.Parallel()
	args := []string{"blah", "blah"}
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile(t)

	want := usage
	habit.RunCLI(tmpFile.Name(), args, &buffer)
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("too many arguments should print usage Message got: %s", got)
	}
}

func TestOptionsButNoArgsShowsUsageHelp(t *testing.T) {
	t.Parallel()
	args := []string{"-f", "daily"}
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile(t)

	want := "Usage"
	habit.RunCLI(tmpFile.Name(), args, &buffer)
	got := buffer.String()

	if !strings.Contains(got, want) {
		t.Errorf("only options and no arguments should print usage Message got: %s", got)
	}
}

func TestOptionServerStartsHTTPServerReturns404NoHabits(t *testing.T) {
	t.Parallel()
	tmpFile := CreateTmpFile(t)
	freePort, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	address := fmt.Sprintf("%s:%d", localHostAddress, freePort)
	args := []string{"-s", address}
	buffer := bytes.Buffer{}

	go habit.RunCLI(tmpFile.Name(), args, &buffer)
	time.Sleep(2 * time.Second)
	resp, err := http.Get("http://" + address)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Want Status %d, got: %d", http.StatusNotFound, resp.StatusCode)
	}

	want := notFound
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Not able to parse response")
	}
	if !strings.Contains(string(got), want) {
		t.Errorf("want response body to be:\n %s \ngot:\n %s", want, got)
	}
}

func TestOptionServerStartsHTTPServerReturns200ExistingHabits(t *testing.T) {
	t.Parallel()
	tmpFile := CreateTmpFile(t)
	freePort, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	address := fmt.Sprintf("%s:%d", localHostAddress, freePort)
	args := []string{"-s", address}
	buffer := bytes.Buffer{}

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

	go habit.RunCLI(tmpFile.Name(), args, &buffer)
	time.Sleep(2 * time.Second)
	resp, err := http.Get("http://" + address)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Want Status %d, got: %d", http.StatusNotFound, resp.StatusCode)
	}

	want := streak
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Not able to parse response")
	}
	if !strings.Contains(string(got), want) {
		t.Errorf("want response body to be:\n %s \ngot:\n %s", want, got)
	}
}

func TestExistingHabitShowsStreakMessage(t *testing.T) {
	t.Parallel()
	args := []string{"piano"}
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile(t)

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

	want := "piano"
	got := buffer.String()
	if !strings.Contains(got, want) {
		t.Errorf("existing habit should print streak message. Got:\n  %s", got)
	}
}

func TestNewHabitShowsStreakMessage(t *testing.T) {
	t.Parallel()
	args := []string{"piano"}
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile(t)

	habit.RunCLI(tmpFile.Name(), args, &buffer)

	want := "piano"
	got := buffer.String()
	if !strings.Contains(got, want) {
		t.Errorf("new habit should print streak message. Got:\n  %s", got)
	}
}

func TestNewHabitInvalidFrequencyFails(t *testing.T) {
	t.Parallel()
	args := []string{"-f", "yellow", "piano"}
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile(t)

	habit.RunCLI(tmpFile.Name(), args, &buffer)

	want := "unknown frequency"
	got := buffer.String()
	if !strings.Contains(got, want) {
		t.Errorf("invalid frequency should print error messasge. Got: \n  %s", got)
	}
}

func TestNewWeeklyHabitShowsStreakMessage(t *testing.T) {
	t.Parallel()
	args := []string{"-f", "weekly", "piano"}
	buffer := bytes.Buffer{}
	tmpFile := CreateTmpFile(t)

	habit.RunCLI(tmpFile.Name(), args, &buffer)

	want := "piano"
	got := buffer.String()
	if !strings.Contains(got, want) {
		t.Errorf("new weekly habit should print streak message. Got:\n  %s", got)
	}
}
