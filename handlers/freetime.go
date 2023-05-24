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

func ApiFreeTimes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["staff_id"]

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var freeTimes []models.FreeTime
	db.DB.Where("staff_id", param).Where("date_regarding > ? AND date_regarding < ?", startDate, endDate).Find(&freeTimes)

	json, err := json.Marshal(freeTimes)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiFreeTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var freeTime models.FreeTime
	db.DB.Where("id", param).Find(&freeTime)

	json, err := json.Marshal(freeTime)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiFreeTimeCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data models.FreeTime
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
