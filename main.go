package main

import (
	"github.com/araquach/apiAuth/routes"
	apiTime "github.com/araquach/apiTime/routes"
	db "github.com/araquach/dbService"
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
	apiTime.TimeRouter()

	log.Printf("Starting server on %s", port)

	http.ListenAndServe(":"+port, &routes.R)
}
