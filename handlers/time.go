package handlers

import (
	"encoding/json"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ApiTimeDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["staff_id"]

	var timeInfo models.Time

	db.DB.Where("staff_id", param).First(&timeInfo)

	json, err := json.Marshal(timeInfo)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiTimeDash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["staff_id"]

	var timeDash models.TimeDash

	const sql = ``

	db.DB.Raw(sql, id).Scan(&timeDash)

	json, err := json.Marshal(timeDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
