package habit_test

import (
	"bytes"
	"fmt"
	"github.com/crmejia/habit"
	"github.com/phayes/freeport"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
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

func TestRunCLIAllShowsCorrectlyOnChangedDir(t *testing.T) {
	t.Parallel()
	args := []string{"-d", "testdata/", "-s", "file", "all"}
	buffer := bytes.Buffer{}
	habit.RunCLI(args, &buffer)

	want := "Habits:"
	got := buffer.String()
	if !strings.Contains(got, want) {
		// if habits exist in store, a summary of all habits should be displayed
		t.Errorf("habit all should print a summary of all habits got:  \n  %s", got)
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

func TestRunServerShowsErrorOnWrongArgs(t *testing.T) {
	t.Parallel()
	const noAddressError = "no address provided"
	const tooManyArgsError = "too many args provided"
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{name: "nil args", args: nil, want: noAddressError},
		{name: "empty args", args: []string{}, want: noAddressError},
		{name: "too many args", args: []string{"blah", "blah"}, want: tooManyArgsError},
	}

	for _, tc := range testCases {
		got := bytes.Buffer{}
		habit.RunServer(tc.args, &got)
		if !strings.Contains(got.String(), tc.want) {
			t.Errorf("%s should fail with %s, got: %s", tc.name, noAddressError, got.String())
		}
	}
}
func TestRunServerStartsServer(t *testing.T) {
	t.Parallel()
	freePort, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	address := fmt.Sprintf("%s:%d", localHostAddress, freePort)
	args := []string{address}
	output := bytes.Buffer{}
	go habit.RunServer(args, &output)

	address = "http://" + address + "?habit=piano"
	resp, err := retryHttpGet(address)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Want Status %d, got: %d", http.StatusNotFound, resp.StatusCode)
	}

	want := "piano"
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Not able to parse response")
	}
	if !strings.Contains(string(got), want) {
		t.Errorf("want response body to be:\n %s \ngot:\n %s", want, got)
	}
}
func retryHttpGet(address string) (*http.Response, error) {
	resp, err := http.Get(address)
	for err != nil {
		switch {
		case strings.Contains(err.Error(), "connection refused"):
			time.Sleep(5 * time.Millisecond)
			resp, err = http.Get(address)
		default:
			return resp, err
		}
	}
	return resp, nil
}
