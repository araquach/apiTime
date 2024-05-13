package routes

import "github.com/gorilla/mux"

func TimeRouter() *mux.Router {
	r := mux.NewRouter()

	// Main routes
	time(r)
	freeTime(r)
	holiday(r)
	lieu(r)
	sick(r)

	//Admin routes
	holidayAdmin(r)

	return r
}
