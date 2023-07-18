package handlers

import (
	"encoding/json"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ApiTimeInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["staff_id"]

	var timeInfo models.Time

	db.DB.Where("staff_id", param).Find(timeInfo)

	json, err := json.Marshal(timeInfo)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
