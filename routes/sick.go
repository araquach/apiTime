package routes

import (
	"github.com/araquach/apiAuth/helpers"
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func sick(r *mux.Router) {
	s := r.PathPrefix("/api/time/sick").Subrouter()

	s.HandleFunc("/dash/{staff_id}", handlers.ApiSickDash).Methods("GET")
	s.HandleFunc("/all/{staff_id}", handlers.ApiSickDays).Methods("GET")
	s.HandleFunc("/sick/{id}", handlers.ApiSickDay).Methods("GET")
	s.HandleFunc("/sick/create", helpers.TokenVerifyMiddleWare(handlers.ApiSickDayCreate)).Methods("POST")
}
