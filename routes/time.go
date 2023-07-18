package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func time(r *mux.Router) {
	s := r.PathPrefix("/api/time").Subrouter()

	s.HandleFunc("/time-info/{staff_id}", handlers.ApiTimeInfo).Methods("GET")
}
