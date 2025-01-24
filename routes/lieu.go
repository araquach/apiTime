package routes

import (
	"github.com/araquach/apiTime/handlers"
	"github.com/gorilla/mux"
)

func lieu(r *mux.Router) {
	s := r.PathPrefix("/api/time/lieu").Subrouter()
	s.HandleFunc("/dash/{staff_id}", handlers.ApiLieuDash).Methods("GET")
	s.HandleFunc("/all/{staff_id}/{year}", handlers.ApiLieuHours).Methods("GET")
	s.HandleFunc("/{id}", handlers.ApiLieuHour).Methods("GET")
	s.HandleFunc("/create", handlers.ApiLieuHourCreate).Methods("POST")
	s.HandleFunc("/update/{id}", handlers.ApiLieuHourUpdate).Methods("PUT")
}

func lieuAdmin(r *mux.Router) {
	s := r.PathPrefix("/api/time/admin/lieu").Subrouter()
	s.HandleFunc("/dash", handlers.ApiLieuAdminDash).Methods("GET")
	s.HandleFunc("/all/pending", handlers.ApiAdminLieuPending).Methods("GET")
	s.HandleFunc("/all/{staff_id}", handlers.ApiAdminLieuHours).Methods("GET")
	s.HandleFunc("/approve/{id}", handlers.ApiLieuApprove).Methods("PUT")
}
