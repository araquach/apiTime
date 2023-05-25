package main

import (
	"github.com/araquach/apiTime/routes"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db.DBInit(dsn)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Load API Routes
	timeRouter := routes.TimeRouter()
	mainRouter := mux.NewRouter()

	mainRouter.PathPrefix("/api/time").Handler(timeRouter)

	log.Printf("Starting server on %s", port)

	http.ListenAndServe(":"+port, mainRouter)
}
