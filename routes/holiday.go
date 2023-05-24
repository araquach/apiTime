package routes

import (
	"github.com/araquach/apiAuth/helpers"
	"github.com/araquach/apiTime/handlers"
)

func holiday() {
	s := R.PathPrefix("/api/time").Subrouter()

	s.HandleFunc("/api/holidays/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiHolidays)).Methods("GET")
	s.HandleFunc("/api/holiday/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiHoliday)).Methods("GET")
	s.HandleFunc("/api/holiday", helpers.TokenVerifyMiddleWare(handlers.ApiHolidayCreate)).Methods("POST")
}