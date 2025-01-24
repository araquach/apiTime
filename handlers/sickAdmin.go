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

func ApiSickAdminDash(w http.ResponseWriter, r *http.Request) {

	var sickAdminDash models.SickAdminDash

	const sql = `SELECT COUNT(*) AS pending FROM sicks WHERE deducted = 0;`

	db.DB.Raw(sql).Scan(&sickAdminDash)

	json, err := json.Marshal(sickAdminDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiAdminSickHours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["staff_id"]

	type Result struct {
		models.Sick
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var results []Result

	db.DB.Table("sicks").Where("sicks.staff_id", id).Select("sicks.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on team_members.staff_id = sicks.staff_id").
		Where("sicks.date_from > ? AND sicks.date_from < ?", startDate, endDate).
		Scan(&results)

	json, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiAdminSickPending(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Result struct {
		models.Sick
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var results []Result

	db.DB.Table("sicks").Select("sicks.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on team_members.staff_id = sicks.staff_id").
		Where("sicks.date_from > ? AND sicks.date_from < ?", startDate, endDate).
		Where("sicks.deducted = ?", 0).
		Scan(&results)

	json, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiSickDayCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var sick models.Sick
	err := decoder.Decode(&sick)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	result := db.DB.Create(&sick)
	if result.Error != nil {
		http.Error(w, "Failed to create record: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func ApiSickDayUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var sick models.Sick

	err := json.NewDecoder(r.Body).Decode(&sick)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()

	var originalSick models.Sick
	res := tx.First(&originalSick, id)
	if res.Error != nil {
		log.Printf("Error finding original sick day: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.Model(&models.Sick{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Description": sick.Description,
		"DateFrom":    sick.DateFrom,
		"DateTo":      sick.DateTo,
		"Hours":       sick.Hours,
	})
	if res.Error != nil {
		log.Printf("Error updating sick: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(sick)
}

func ApiSickDayDeduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var sick models.Sick

	err := json.NewDecoder(r.Body).Decode(&sick)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()

	var originalSick models.Sick
	res := tx.First(&originalSick, id)
	if res.Error != nil {
		log.Printf("Error finding original sick day: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.Model(&models.Sick{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Deducted": sick.Deducted,
	})
	if res.Error != nil {
		log.Printf("Error updating Sick Day: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(sick)
}
