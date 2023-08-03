package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func freeTime(r *mux.Router) {
	s := r.PathPrefix("/api/time").Subrouter()

	//s.HandleFunc("/free-times/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiFreeTimes)).Methods("GET")
	s.HandleFunc("/free-times/{staff_id}", handlers.ApiFreeTimes).Methods("GET")
	//s.HandleFunc("/free-time/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiFreeTime)).Methods("GET")
	s.HandleFunc("/free-time/{id}", handlers.ApiFreeTime).Methods("GET")
	//s.HandleFunc("/free-time-create", helpers.TokenVerifyMiddleWare(handlers.ApiFreeTimeCreate)).Methods("POST")
	s.HandleFunc("/free-time-create", handlers.ApiFreeTimeCreate).Methods("POST")
	s.HandleFunc("/free-time-update/{id}", handlers.ApiFreeTimeUpdate).Methods("PUT")
}
