package routes

import (
	"github.com/araquach/apiAuth/helpers"
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func sick(r *mux.Router) {
	s := r.PathPrefix("/api/time").Subrouter()

	s.HandleFunc("/sick-days/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiSickDays)).Methods("GET")
	s.HandleFunc("/sick-day/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiSickDay)).Methods("GET")
	s.HandleFunc("/sick-day", helpers.TokenVerifyMiddleWare(handlers.ApiSickDayCreate)).Methods("POST")
}
