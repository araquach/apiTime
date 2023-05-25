package routes

import "github.com/gorilla/mux"

func TimeRouter() *mux.Router {
	r := mux.NewRouter()

	freeTime(r)
	holiday(r)
	lieu(r)
	sick(r)

	return r
}
