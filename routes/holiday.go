package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func holiday(r *mux.Router) {
	s := r.PathPrefix("/api/time/holiday").Subrouter()
	s.HandleFunc("/dash/{staff_id}", handlers.ApiHolidayDash).Methods("GET")
	//s.HandleFunc("/holidays/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiHolidays)).Methods("GET")
	s.HandleFunc("/all/{staff_id}", handlers.ApiHolidays).Methods("GET")
	//s.HandleFunc("/holiday/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiHoliday)).Methods("GET")
	s.HandleFunc("/{id}", handlers.ApiHoliday).Methods("GET")
	//s.HandleFunc("/holiday-create", helpers.TokenVerifyMiddleWare(handlers.ApiHolidayCreate)).Methods("POST")
	s.HandleFunc("/create", handlers.ApiHolidayCreate).Methods("POST")
	s.HandleFunc("/update/{id}", handlers.ApiHolidayUpdate).Methods("PUT")
}

func holidayAdmin(r *mux.Router) {
	s := r.PathPrefix("/api/time/admin/holiday").Subrouter()
	s.HandleFunc("/approve/{id}", handlers.ApiHolidayApprove).Methods("PUT")
}
