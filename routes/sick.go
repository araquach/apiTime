package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func sick(r *mux.Router) {
	s := r.PathPrefix("/api/time/sick").Subrouter()
	s.HandleFunc("/dash/{staff_id}", handlers.ApiSickDash).Methods("GET")
	s.HandleFunc("/all/{staff_id}/{year}", handlers.ApiSickDays).Methods("GET")
	s.HandleFunc("/{id}", handlers.ApiSickDay).Methods("GET")
}

func sickAdmin(r *mux.Router) {
	s := r.PathPrefix("/api/time/admin/sick").Subrouter()
	s.HandleFunc("/dash", handlers.ApiSickAdminDash).Methods("GET")
	s.HandleFunc("/create", handlers.ApiSickDayCreate).Methods("POST")
	s.HandleFunc("/all/pending", handlers.ApiAdminSickPending).Methods("GET")
	s.HandleFunc("/update/{id}", handlers.ApiSickDayUpdate).Methods("PUT")
	s.HandleFunc("/all/{staff_id}", handlers.ApiAdminSickHours).Methods("GET")
	s.HandleFunc("/deduct/{id}", handlers.ApiSickDayDeduct).Methods("PUT")
}
