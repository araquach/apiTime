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
	db.DB.Where("staff_id", param).Where("request_date > ? AND request_date < ?", startDate, endDate).Find(&freeTimes)

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

	res := tx.Where("staff_id", freeTime.StaffId).First(&time)
	if res.Error != nil {
		// Handle error, e.g., log it and return
		log.Printf("Error finding booking: %v", res.Error)
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

	res = tx.Model(&models.FreeTime{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Description": freeTime.Description,
		"RequestDate": freeTime.RequestDate,
		"Hours":       freeTime.Hours,
	})
	if res.Error != nil {
		log.Printf("Error finding time entry: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(freeTime)
}

func ApiFreeTimeDash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["staff_id"]

	var freeTimeDash models.FreeTimeDash

	const sql = `SELECT
    MAX(times.free_time_ent) AS "entitlement",
    COALESCE(MAX(times.free_time_ent) - SUM(CASE WHEN free_times.approved = 1 THEN free_times.hours ELSE 0 END), MAX(times.free_time_ent)) AS "remaining",
    COALESCE(MAX(times.free_time_ent) - SUM(CASE WHEN free_times.approved = 0 THEN free_times.hours ELSE 0 END), MAX(times.free_time_ent)) AS "remaining_pending",
    COALESCE(SUM(CASE WHEN free_times.approved = 1 THEN free_times.hours ELSE 0 END), 0) AS "used",
    COALESCE(SUM(CASE WHEN free_times.approved = 0 THEN free_times.hours ELSE 0 END), 0) AS "pending"
FROM
    times
        LEFT JOIN free_times ON times.staff_id = free_times.staff_id
WHERE
    times.staff_id = ?
GROUP BY
    times.free_time_ent`

	db.DB.Raw(sql, id).Scan(&freeTimeDash)

	json, err := json.Marshal(freeTimeDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
