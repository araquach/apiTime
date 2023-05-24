package routes

import (
	"github.com/araquach/apiAuth/helpers"
	"github.com/araquach/apiTime/handlers"
)

func sick() {
	s := R.PathPrefix("/api/time").Subrouter()

	s.HandleFunc("/sick-days/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiSickDays)).Methods("GET")
	s.HandleFunc("/sick-day/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiSickDay)).Methods("GET")
	s.HandleFunc("/sick-day", helpers.TokenVerifyMiddleWare(handlers.ApiSickDayCreate)).Methods("POST")
}
