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

func ApiFreeTimeAdminDash(w http.ResponseWriter, r *http.Request) {

	var freeTimeAdminDash models.FreeTimeAdminDash

	const sql = `SELECT COUNT(*) AS pending FROM free_times WHERE approved = 0;`

	db.DB.Raw(sql).Scan(&freeTimeAdminDash)

	json, err := json.Marshal(freeTimeAdminDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiAdminFreeTimePending(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Result struct {
		models.FreeTime
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var results []Result

	db.DB.Table("free_times").Select("free_times.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on team_members.staff_id = free_times.staff_id").
		Where("free_times.request_date > ? AND free_times.request_date < ?", startDate, endDate).
		Where("free_times.approved = ?", 0).
		Scan(&results)

	json, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiAdminFreeTimeHours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["staff_id"]

	type Result struct {
		models.FreeTime
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var results []Result

	db.DB.Table("free_times").Where("free_times.staff_id", id).Select("free_times.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on team_members.staff_id = free_times.staff_id").
		Where("free_times.request_date > ? AND free_times.request_date < ?", startDate, endDate).
		Scan(&results)

	json, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiFreeTimeApprove(w http.ResponseWriter, r *http.Request) {
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
		"Approved": freeTime.Approved,
	})
	if res.Error != nil {
		log.Printf("Error updating free time: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(freeTime)
}
