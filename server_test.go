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
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())

	server := habit.NewServer(tmpFile.Name(), address)

	//TODO should the tracker be exposed?
	if server.Tracker == nil {
		t.Errorf("Tracker should not be nil")
	}
}

func TestNewServerWithNonDefaultAddress(t *testing.T) {
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())
	want := "http://test.net:8080"
	server := habit.NewServer(tmpFile.Name(), want)

	got := server.Server.Addr
	if want != got {
		t.Errorf("want address to be %s, got %s", want, got)
	}

}

func TestAllHabitsHandlerReturnsAllHabits(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/all", nil)
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())
	server := habit.NewServer(tmpFile.Name(), address)
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

func TestHabitHandlerReturnsHabit(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?habit=piano", nil)
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())
	server := habit.NewServer(tmpFile.Name(), address)
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

func TestHabitHandlerWithGibberishReturns400(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?garbage", nil)
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())
	server := habit.NewServer(tmpFile.Name(), address)
	server.HabitHandler(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want status code %d, got %d", http.StatusBadRequest, res.StatusCode)
	}

}

func TestRouting(t *testing.T) {
	//todo Test index(/) returns all habits
}

const address = "127.0.0.1:8080"
