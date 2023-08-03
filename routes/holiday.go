package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func holiday(r *mux.Router) {
	s := r.PathPrefix("/api/time").Subrouter()

	//s.HandleFunc("/holidays/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiHolidays)).Methods("GET")
	s.HandleFunc("/holidays/{staff_id}", handlers.ApiHolidays).Methods("GET")
	//s.HandleFunc("/holiday/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiHoliday)).Methods("GET")
	s.HandleFunc("/holiday/{id}", handlers.ApiHoliday).Methods("GET")
	//s.HandleFunc("/holiday-create", helpers.TokenVerifyMiddleWare(handlers.ApiHolidayCreate)).Methods("POST")
	s.HandleFunc("/holiday-create", handlers.ApiHolidayCreate).Methods("POST")
	s.HandleFunc("/holiday-update/{id}", handlers.ApiHolidayUpdate).Methods("PUT")
}
