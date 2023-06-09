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

func ApiHolidays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["staff_id"]

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var holidays []models.Holiday
	db.DB.Where("staff_id", param).Where("request_date_from > ? AND request_date_from < ?", startDate, endDate).Find(&holidays)

	json, err := json.Marshal(holidays)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiHoliday(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var holiday models.Holiday
	db.DB.Where("id", param).Find(&holiday)

	json, err := json.Marshal(holiday)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiHolidayCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data models.Holiday
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
