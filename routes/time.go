package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func time(r *mux.Router) {
	s := r.PathPrefix("/api/time").Subrouter()
	s.HandleFunc("/dash/{staff_id}", handlers.ApiTimeDash).Methods("GET")
	s.HandleFunc("/time-details/{staff_id}", handlers.ApiTimeDetails).Methods("GET")
}

func timeAdmin(r *mux.Router) {
	s := r.PathPrefix("/api/time/admin").Subrouter()
	s.HandleFunc("/dash/{salon_id}", handlers.ApiTimeAdminDash).Methods("GET")
}
