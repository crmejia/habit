package habit_test

import (
	"habit"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	localHostAddress = "127.0.0.1"
	notFound         = "not found"
)

func TestNewHttpServer(t *testing.T) {
	t.Parallel()
	tmpFile := CreateTmpFile(t)

	store := habit.NewFileStore(tmpFile.Name())
	server := habit.NewServer(store, localHostAddress)

	if server.Tracker == nil {
		t.Errorf("Tracker should not be nil")
	}
}

func TestNewServerWithNonDefaultAddress(t *testing.T) {
	t.Parallel()
	tmpFile := CreateTmpFile(t)
	want := "http://test.net:8080"
	store := habit.NewFileStore(tmpFile.Name())
	server := habit.NewServer(store, want)

	got := server.Server.Addr
	if want != got {
		t.Errorf("want localHostAddress to be %s, got %s", want, got)
	}

}

func TestHabitHandlerReturnsHabit(t *testing.T) {
	t.Parallel()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?habit=piano", nil)
	tmpFile := CreateTmpFile(t)
	store := habit.NewFileStore(tmpFile.Name())
	server := habit.NewServer(store, localHostAddress)
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
	t.Parallel()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?garbage", nil)
	tmpFile := CreateTmpFile(t)
	store := habit.NewFileStore(tmpFile.Name())
	server := habit.NewServer(store, localHostAddress)
	server.HabitHandler(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want status code %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestServer_HabitHandleInterval(t *testing.T) {
	t.Parallel()
	tmpFile := CreateTmpFile(t)
	store := habit.NewFileStore(tmpFile.Name())
	server := habit.NewServer(store, localHostAddress)

	testCases := []struct {
		name   string
		target string
		want   int
	}{
		{"BadRequest on invalid interval", "/?habit=piano&interval=wrong", http.StatusBadRequest},
		{"OK on valid interval", "/?habit=piano&interval=weekly", http.StatusOK},
	}

	for _, tc := range testCases {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, tc.target, nil)
		server.HabitHandler(recorder, req)
		res := recorder.Result()
		got := res.StatusCode
		if tc.want != got {
			t.Errorf("want %s, got: %d", tc.name, got)
		}
	}
}

func TestRouting(t *testing.T) {
	t.Parallel()
	tmpFile := CreateTmpFile(t)
	store := habit.NewFileStore(tmpFile.Name())
	habitServer := habit.NewServer(store, localHostAddress)
	testServer := httptest.NewServer(habitServer.Handler())
	defer testServer.Close()

	testCases := []struct {
		path           string
		wantStatusCode int
	}{
		{path: "/", wantStatusCode: http.StatusOK},
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
