package habit_test

import (
	"fmt"
	"github.com/crmejia/habit"
	"github.com/phayes/freeport"
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
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}
	server, err := habit.NewServer(&controller, localHostAddress)
	if err != nil {
		t.Fatal(err)
	}

	if server.Server == nil {
		t.Errorf("Tracker should not be nil")
	}
}

func TestNewServerWithNonDefaultAddress(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}
	want := "http://test.net:8080"
	server, err := habit.NewServer(&controller, want)
	if err != nil {
		t.Fatal(err)
	}

	got := server.Server.Addr
	if want != got {
		t.Errorf("want localHostAddress to be %s, got %s", want, got)
	}

}
func TestNewServerReturnsErrorOnEmptyAddress(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}
	_, err = habit.NewServer(&controller, "")
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
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}
	server, err := habit.NewServer(&controller, localHostAddress)
	if err != nil {
		t.Fatal(err)
	}

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
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}

	server, err := habit.NewServer(&controller, localHostAddress)
	if err != nil {
		t.Fatal(err)
	}
	handler := server.HandleIndex()
	handler(recorder, req)

	res := recorder.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want status code %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
	err = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServer_HabitHandleFrequency(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}
	server, err := habit.NewServer(&controller, localHostAddress)
	if err != nil {
		t.Fatal(err)
	}

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
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}
	server, err := habit.NewServer(&controller, localHostAddress)
	if err != nil {
		t.Fatal(err)
	}

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
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}
	habitServer, err := habit.NewServer(&controller, localHostAddress)
	if err != nil {
		t.Fatal(err)
	}
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
		err = res.Body.Close()
		if err != nil {
			t.Fatal(err)
		} //no defer at it might leak ;)
	}
}

func TestServer_RunReturnsBadRequest(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}

	freePort, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	address := fmt.Sprintf("%s:%d", localHostAddress, freePort)
	server, err := habit.NewServer(&controller, address)
	if err != nil {
		t.Fatal(err)
	}
	go server.Run()
	resp, err := retryHttpGet("http://" + address)
	if err != nil {
		t.Fatal(err)

	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Want Status %d, got: %d", http.StatusNotFound, resp.StatusCode)
	}

	want := "cannot parse querystring"
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Not able to parse response")
	}
	if !strings.Contains(string(got), want) {
		t.Errorf("want response body to be:\n %s \ngot:\n %s", want, got)
	}
}
func TestServer_RunReturnsHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, err := habit.NewController(&store)
	if err != nil {
		t.Fatal(err)
	}

	freePort, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	address := fmt.Sprintf("%s:%d", localHostAddress, freePort)
	server, err := habit.NewServer(&controller, address)
	if err != nil {
		t.Fatal(err)
	}
	go server.Run()
	resp, err := retryHttpGet("http://" + address + "?habit=piano")
	for err != nil {
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
