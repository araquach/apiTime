package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func holiday(r *mux.Router) {
	s := r.PathPrefix("/api/time/holiday").Subrouter()
	s.HandleFunc("/dash/{staff_id}", handlers.ApiHolidayDash).Methods("GET")
	s.HandleFunc("/all/{staff_id}", handlers.ApiHolidays).Methods("GET")
	s.HandleFunc("/{id}", handlers.ApiHoliday).Methods("GET")
	s.HandleFunc("/create", handlers.ApiHolidayCreate).Methods("POST")
	s.HandleFunc("/update/{id}", handlers.ApiHolidayUpdate).Methods("PUT")
}

func holidayAdmin(r *mux.Router) {
	s := r.PathPrefix("/api/time/admin/holiday").Subrouter()
	s.HandleFunc("/dash", handlers.ApiHolidayAdminDash).Methods("GET")
	s.HandleFunc("/all/pending", handlers.ApiAdminHolidaysPending).Methods("GET")
	s.HandleFunc("/all/{staff_id}", handlers.ApiAdminHolidays).Methods("GET")
	s.HandleFunc("/approve/{id}", handlers.ApiHolidayApprove).Methods("PUT")
}
