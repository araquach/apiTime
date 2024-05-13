package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func lieu(r *mux.Router) {
	s := r.PathPrefix("/api/time/lieu").Subrouter()
	s.HandleFunc("/dash/{staff_id}", handlers.ApiLieuDash).Methods("GET")
	//s.HandleFunc("/lieu-hours/{staff_id}", helpers.TokenVerifyMiddleWare(handlers.ApiLieuHours)).Methods("GET")
	s.HandleFunc("/all/{staff_id}", handlers.ApiLieuHours).Methods("GET")
	//s.HandleFunc("/lieu-hour/{id}", helpers.TokenVerifyMiddleWare(handlers.ApiLieuHour)).Methods("GET")
	s.HandleFunc("/{id}", handlers.ApiLieuHour).Methods("GET")
	//s.HandleFunc("/lieu-hour-create", helpers.TokenVerifyMiddleWare(handlers.ApiLieuHourCreate)).Methods("POST")
	s.HandleFunc("/create", handlers.ApiLieuHourCreate).Methods("POST")
	s.HandleFunc("/update/{id}", handlers.ApiLieuHourUpdate).Methods("PUT")
}
