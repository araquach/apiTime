package routes

import (
	"github.com/araquach/apiAuth/helpers"
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func freeTime(r *mux.Router) {
	s := r.PathPrefix("/api/time").Subrouter()

	s.HandleFunc("/free-times/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiFreeTimes)).Methods("GET")
	s.HandleFunc("/free-time/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiFreeTime)).Methods("GET")
	s.HandleFunc("/free-time", helpers.TokenVerifyMiddleWare(handlers.ApiFreeTimeCreate)).Methods("POST")
}
