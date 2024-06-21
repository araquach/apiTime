package routes

import "github.com/gorilla/mux"

func TimeRouter() *mux.Router {
	r := mux.NewRouter()

	// Main routes
	time(r)
	holiday(r)
	lieu(r)
	freeTime(r)
	sick(r)

	//Admin routes
	timeAdmin(r)
	holidayAdmin(r)
	lieuAdmin(r)
	freeTimeAdmin(r)
	sickAdmin(r)

	return r
}
