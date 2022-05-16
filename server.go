package habit

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	defaultTCPAddress = "127.0.0.1:8080"
)

type server struct {
	Server  *http.Server
	Tracker *Tracker
	Store   Store
}

//NewServer returns a new server
func NewServer(store Store, address string) *server {
	tracker := NewTracker(store)
	server := server{
		Server: &http.Server{
			Addr: address},
		Tracker: &tracker,
		Store:   store,
	}
	return &server
}

//Run listens and serves http
func (server *server) Run() {
	router := server.Router()
	server.Server.Handler = router

	err := server.Server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Println(err)
	}
}

//Router returns a http.Handler with the appropriate routes
func (server *server) Router() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", server.HabitHandler)

	return router
}

//HabitHandler Handler that servers habits.
func (server *server) HabitHandler(w http.ResponseWriter, r *http.Request) {
	//parsing querystring
	habitName := r.FormValue("habit")
	if r.RequestURI == "/all" || r.RequestURI == "/" {
		if len(*server.Tracker) > 0 {
			fmt.Fprint(w, AllHabits(server.Store))
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if habitName == "" || r.URL.Path != "/" {
		http.Error(w, "cannot parse querystring", http.StatusBadRequest)
		return
	}

	intervalString := r.FormValue("interval")
	var interval time.Duration
	if intervalString == "" || intervalString == "daily" {
		interval = DailyInterval
	} else if intervalString == "weekly" {
		interval = WeeklyInterval
	} else {
		http.Error(w, "invalid interval", http.StatusBadRequest)
		return
	}

	habit, ok := server.Tracker.FetchHabit(habitName)
	if !ok {
		habit = &Habit{Name: habitName, Interval: interval}

		// this error is hard to test as the conditions that trigger cannot be met:
		// - trying to create an existing habit(which FetchHabit(ln67) covers
		// - invalid interval which ln59-65 cover
		err := server.Tracker.CreateHabit(habit)
		if err != nil {
			http.Error(w, "not able to create habit", http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprint(w, habit)
}
