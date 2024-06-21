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

func ApiHolidayAdminDash(w http.ResponseWriter, r *http.Request) {

	var holidayAdminDash models.HolidayAdminDash

	const sql = `SELECT COUNT(*) AS pending FROM holidays WHERE approved = 0;`

	db.DB.Raw(sql).Scan(&holidayAdminDash)

	json, err := json.Marshal(holidayAdminDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiAdminHolidaysPending(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Result struct {
		models.Holiday
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var results []Result

	db.DB.Table("holidays").Select("holidays.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on team_members.staff_id = holidays.staff_id").
		Where("holidays.date_from > ? AND holidays.date_from < ?", startDate, endDate).
		Where("holidays.approved = ?", 0).
		Scan(&results)

	json, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiAdminHolidays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["staff_id"]

	type Result struct {
		models.Holiday
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var results []Result

	db.DB.Table("holidays").Where("holidays.staff_id", id).Select("holidays.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on team_members.staff_id = holidays.staff_id").
		Where("holidays.date_from > ? AND holidays.date_from < ?", startDate, endDate).
		Scan(&results)

	json, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiHolidayApprove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var holiday models.Holiday
	// Decode the incoming holiday JSON body
	err := json.NewDecoder(r.Body).Decode(&holiday)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Start a transaction
	tx := db.DB.Begin()
	// Find the original holiday
	var originalHoliday models.Holiday
	res := tx.First(&originalHoliday, id)
	if res.Error != nil {
		log.Printf("Error finding original holiday: %v", res.Error)
		tx.Rollback()
		return
	}
	// Update the holiday in the database
	res = tx.Model(&models.Holiday{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Approved": holiday.Approved,
	})
	if res.Error != nil {
		log.Printf("Error updating holiday: %v", res.Error)
		tx.Rollback()
		return
	}
	// If everything went well, commit the transaction
	tx.Commit()
	// Return the updated holiday
	json.NewEncoder(w).Encode(holiday)
}
