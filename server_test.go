package habit_test

import (
	"habit"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNewHttpServer(t *testing.T) {
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())
	server := habit.NewServer(tmpFile.Name())

	//TODO should the tracker be exposed?
	if server.Tracker == nil {
		t.Errorf("Tracker should not be nil")
	}

}

func TestGetReturnsAllHabits(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/all", nil)
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())
	server := habit.NewServer(tmpFile.Name())
	server.AllHabitsHandler(recorder, request)
	res := recorder.Result()
	defer res.Body.Close()

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("couldn't read response:%v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got: %d", res.StatusCode)
	}

	want := "Habits:"
	if !strings.Contains(string(got), want) {
		t.Errorf("wanted a list of habits, got:\n %s", got)
	}
}

func TestGetReturnsHabit(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?habit=piano", nil)
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())
	server := habit.NewServer(tmpFile.Name())
	server.Tracker = &habit.Tracker{
		"reading": &habit.Habit{
			Name: "reading",
		},
	}

	server.HabitHandler(recorder, req)
	res := recorder.Result()
	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("couldn't read response:%v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got: %d", res.StatusCode)
	}
	want := "piano"
	if !strings.Contains(string(got), want) {
		t.Errorf("expected habit '%s', got:\n%s", want, got)
	}
}

func TestGetWithGibberishReturns400(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?garbage", nil)
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())
	server := habit.NewServer(tmpFile.Name())
	server.HabitHandler(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want status code %d, got %d", http.StatusBadRequest, res.StatusCode)
	}

}
