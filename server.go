package habit

import (
	"fmt"
	"log"
	"net/http"
)

type server struct {
	Server  *http.Server
	Tracker *Tracker
}

func NewServer(filename string, address string) *server {
	//todo inject tracker
	tracker := NewTracker(filename)
	server := server{
		Server:  &http.Server{Addr: address},
		Tracker: &tracker,
	}
	return &server
}

func (server *server) Run() {
	http.HandleFunc("/", server.HabitHandler)
	http.HandleFunc("/all", server.AllHabitsHandler)

	//http.ListenAndServe("localhost:8080", nil)
	err := server.Server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Println(err)
	}
}

func (server *server) HabitHandler(w http.ResponseWriter, r *http.Request) {
	habitName := r.FormValue("habit")

	if habitName == "" {
		http.Error(w, "cannot parse querystring", http.StatusBadRequest)
		return
	}

	habit, ok := server.Tracker.FetchHabit(habitName)
	if !ok {
		habit = &Habit{Name: habitName, Interval: WeeklyInterval}
		//TODO habit interval
		server.Tracker.CreateHabit(habit)
		//if err != nil {
		//	//TODO return error
		//}
	}
	fmt.Fprint(w, habit)
}

func (server *server) AllHabitsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, server.Tracker.AllHabits())
}

const (
	defaultTCPAddress = "127.0.0.1:8080"
)
