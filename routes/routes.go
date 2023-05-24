package routes

import "github.com/gorilla/mux"

var R mux.Router

func TimeRouter() {
	freeTime()
	holiday()
	lieu()
	sick()

	return
}
