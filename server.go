package habit

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type server struct {
	*http.Server
	controller *Controller
}

//NewServer returns a new server
func NewServer(controller *Controller, address string) (*server, error) {
	if controller == nil {
		return nil, errors.New("controller cannot be nil")
	}
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}

	server := server{
		Server: &http.Server{
			Addr: address},
		controller: controller,
	}
	return &server, nil
}

//Run listens and serves http
func (server *server) Run() {
	router := server.Routes()
	server.Handler = router

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Println(err)
	}
}

//Routes returns a http.Handler with the appropriate routes
func (server *server) Routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", server.HandleIndex())
	router.HandleFunc("/all", server.HandleAll())

	return router
}

//HandleIndex handler that servers habits.
func (server *server) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parsing querystring
		habitName := r.FormValue("habit")
		if habitName == "" || r.URL.Path != "/" {
			http.Error(w, "cannot parse querystring", http.StatusBadRequest)
			return
		}

		frequency := r.FormValue("frequency")
		if frequency == "" {
			frequency = "daily" //default frequency
		}

		inputHabit, err := parseHabit(habitName, frequency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h, err := server.controller.Handle(inputHabit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		fmt.Fprint(w, h)
	}
}

//HandleAll handler that serves /all
func (server *server) HandleAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allHabits := server.controller.GetAllHabits()
		fmt.Fprint(w, allHabits)
	}
}
