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

func ApiLieuHours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["staff_id"]

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var lieuHours []models.Lieu
	db.DB.Where("staff_id", param).Where("date_regarding > ? AND date_regarding < ?", startDate, endDate).Find(&lieuHours)

	json, err := json.Marshal(lieuHours)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiLieuHour(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var lieuHour models.Lieu
	db.DB.Where("id", param).Find(&lieuHour)

	json, err := json.Marshal(lieuHour)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiLieuHourCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data models.Lieu
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
