package routes

import (
	"github.com/araquach/apiAuth/helpers"
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func lieu(r *mux.Router) {
	s := r.PathPrefix("/api/time").Subrouter()

	s.HandleFunc("/lieu-hours/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiLieuHours)).Methods("GET")
	s.HandleFunc("/lieu-hour/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiLieuHour)).Methods("GET")
	s.HandleFunc("/lieu-hour", helpers.TokenVerifyMiddleWare(handlers.ApiLieuHourCreate)).Methods("POST")
}
