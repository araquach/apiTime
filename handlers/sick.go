package handlers

import (
	"encoding/json"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"github.com/jinzhu/now"
	"log"
	"net/http"
)

func ApiSickDays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["staff_id"]

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var sickDays []models.Sick
	db.DB.Where("staff_id", param).Where("sick_from > ? AND sick_from < ?", startDate, endDate).Find(&sickDays)

	json, err := json.Marshal(sickDays)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiSickDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var sickDay models.Sick
	db.DB.Where("id", param).Find(&sickDay)

	json, err := json.Marshal(sickDay)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiSickDayCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data models.Sick
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	db.DB.Create(&data)
	if err != nil {
		log.Fatal(err)
	}
	return
}
