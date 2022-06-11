package habit_test

import (
	"github.com/crmejia/habit"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	localHostAddress = "127.0.0.1"
)

func TestNewHttpServer(t *testing.T) {
	t.Parallel()

	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	server, _ := habit.NewServer(&controller, localHostAddress)

	if server.Server == nil {
		t.Errorf("Tracker should not be nil")
	}
}

func TestNewServerWithNonDefaultAddress(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	want := "http://test.net:8080"
	server, _ := habit.NewServer(&controller, want)

	got := server.Server.Addr
	if want != got {
		t.Errorf("want localHostAddress to be %s, got %s", want, got)
	}

}
func TestNewServerReturnsErrorOnEmptyAddress(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	_, err := habit.NewServer(&controller, "")
	if err == nil {
		t.Error("want NewController to return error on nil store")
	}
}

func TestNewServerReturnsErrorOnNilController(t *testing.T) {
	t.Parallel()
	_, err := habit.NewServer(nil, "http://test.net")
	if err == nil {
		t.Error("want NewController to return error on nil store")
	}
}

func TestHabitHandlerReturnsHabit(t *testing.T) {
	t.Parallel()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?habit=piano", nil)

	store := habit.OpenMemoryStore()
	store.Habits["piano"] = &habit.Habit{
		Name: "piano",
	}
	controller, _ := habit.NewController(&store)
	server, _ := habit.NewServer(&controller, localHostAddress)

	handler := server.HandleIndex()
	handler(recorder, req)
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

func TestHandleIndexWithGibberishReturns400(t *testing.T) {
	t.Parallel()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?garbage", nil)

	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	server, _ := habit.NewServer(&controller, localHostAddress)

	handler := server.HandleIndex()
	handler(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want status code %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestServer_HabitHandleFrequency(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	server, _ := habit.NewServer(&controller, localHostAddress)

	testCases := []struct {
		name   string
		target string
		want   int
	}{
		{"BadRequest on invalid frequency", "/?habit=piano&frequency=wrong", http.StatusBadRequest},
		{"OK on valid frequency", "/?habit=piano&frequency=weekly", http.StatusOK},
	}

	for _, tc := range testCases {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, tc.target, nil)
		handler := server.HandleIndex()
		handler(recorder, req)
		res := recorder.Result()
		got := res.StatusCode
		if tc.want != got {
			t.Errorf("want %s, got: %d", tc.name, got)
		}
	}
}

func TestHandleAllReturnsAllHabits(t *testing.T) {
	t.Parallel()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/all", nil)

	store := habit.OpenMemoryStore()
	store.Habits = map[string]*habit.Habit{
		"piano":   {Name: "piano"},
		"reading": {Name: "reading"},
	}
	controller, _ := habit.NewController(&store)
	server, _ := habit.NewServer(&controller, localHostAddress)

	handler := server.HandleAll()
	handler(recorder, req)
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

func TestRouting(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	store.Habits = map[string]*habit.Habit{
		"piano":   {Name: "piano"},
		"reading": {Name: "reading"},
	}
	controller, _ := habit.NewController(&store)
	habitServer, _ := habit.NewServer(&controller, localHostAddress)
	testServer := httptest.NewServer(habitServer.Routes())
	defer testServer.Close()
	testCases := []struct {
		path           string
		wantStatusCode int
	}{
		{path: "/", wantStatusCode: http.StatusBadRequest},
		{path: "/all", wantStatusCode: http.StatusOK},
		{path: "/?habit=piano", wantStatusCode: http.StatusOK},
		{path: "/piano", wantStatusCode: http.StatusBadRequest},
		{path: "/piano?habit=piano", wantStatusCode: http.StatusBadRequest},
		{path: "/?test=test", wantStatusCode: http.StatusBadRequest},
		{path: "/test?test=test", wantStatusCode: http.StatusBadRequest},
	}

	for _, tc := range testCases {
		res, err := http.Get(testServer.URL + tc.path)
		if err != nil {
			t.Fatalf("could not send http request got error %v", err)
		}
		got := res.StatusCode
		if tc.wantStatusCode != got {
			t.Errorf("want status %d for path:%s, got %d", tc.wantStatusCode, tc.path, got)
		}
		res.Body.Close() //no defer at it might leak ;)
	}
}
