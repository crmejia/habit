package habit

import (
	"fmt"
	"log"
	"net/http"
)

//todo chicken or the egg?
//Should the Tracker have server inside? but since Tracker is a map
//it is not possible. Think about the best away to represent this relationship
type server struct {
	Server  *http.Server
	Tracker *Tracker
}

func NewServer(filename string, address string) *server {
	//todo inject tracker
	tracker := NewTracker(filename)
	server := server{
		Server: &http.Server{
			Addr: address},
		Tracker: &tracker,
	}
	return &server
}

func (server *server) Run() {
	router := server.Handler()
	server.Server.Handler = router

	err := server.Server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Println(err)
	}
}

func (server *server) Handler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", server.HabitHandler)

	return router
}

func (server *server) HabitHandler(w http.ResponseWriter, r *http.Request) {
	habitName := r.FormValue("habit")
	if r.RequestURI == "/all" || r.RequestURI == "/" {
		fmt.Fprint(w, server.Tracker.AllHabits())
		return
	} else if habitName == "" || r.URL.Path != "/" {
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

const (
	defaultTCPAddress = "127.0.0.1:8080"
)
