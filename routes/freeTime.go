package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func freeTime(r *mux.Router) {
	s := r.PathPrefix("/api/time/free-time").Subrouter()
	s.HandleFunc("/dash/{staff_id}", handlers.ApiFreeTimeDash).Methods("GET")
	s.HandleFunc("/all/{staff_id}/{year}", handlers.ApiFreeTimes).Methods("GET")
	s.HandleFunc("/{id}", handlers.ApiFreeTime).Methods("GET")
	s.HandleFunc("/create", handlers.ApiFreeTimeCreate).Methods("POST")
	s.HandleFunc("/update/{id}", handlers.ApiFreeTimeUpdate).Methods("PUT")
}

func freeTimeAdmin(r *mux.Router) {
	s := r.PathPrefix("/api/time/admin/free-time").Subrouter()
	s.HandleFunc("/dash", handlers.ApiFreeTimeAdminDash).Methods("GET")
	s.HandleFunc("/all/pending", handlers.ApiAdminFreeTimePending).Methods("GET")
	s.HandleFunc("/all/{staff_id}", handlers.ApiAdminFreeTimeHours).Methods("GET")
	s.HandleFunc("/approve/{id}", handlers.ApiFreeTimeApprove).Methods("PUT")
}
