package routes

import "github.com/gorilla/mux"

func TimeRouter() *mux.Router {
	r := mux.NewRouter()

	time(r)
	freeTime(r)
	holiday(r)
	lieu(r)
	sick(r)

	return r
}
