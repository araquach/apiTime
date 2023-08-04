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

	var freeTime models.FreeTime
	var time models.Time

	err := decoder.Decode(&freeTime)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()

	res := tx.First(&time)
	if res.Error != nil {
		// Handle error, e.g., log it and return
		log.Printf("Error finding time details: %v", res.Error)
		tx.Rollback()
		return
	}

	time.FreeTime += freeTime.FreeTimeHours

	res = tx.Model(&time).UpdateColumn("free_time", time.FreeTime)
	if res.Error != nil {
		// Handle error, e.g., log it and return
		log.Printf("Error updating free time in times: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.Create(&freeTime)
	if res.Error != nil {
		log.Printf("Error  creating free time: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

func ApiFreeTimeUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var freeTime models.FreeTime
	var time models.Time

	err := json.NewDecoder(r.Body).Decode(&freeTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()

	var originalFreeTime models.FreeTime
	res := tx.First(&originalFreeTime, id)
	if res.Error != nil {
		log.Printf("Error finding original free time: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.First(&time)
	if res.Error != nil {
		log.Printf("Error finding time entry: %v", res.Error)
		tx.Rollback()
		return
	}

	freeTimeDiff := freeTime.FreeTimeHours - originalFreeTime.FreeTimeHours

	time.FreeTimePending += freeTimeDiff

	res = tx.Model(&time).Update("free_time_pending", time.FreeTimePending)
	if res.Error != nil {
		log.Printf("Error updating time details: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.Model(&models.FreeTime{}).Where("id = ?", id).Updates(freeTime)
	if res.Error != nil {
		log.Printf("Error updating free time: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(freeTime)

}
