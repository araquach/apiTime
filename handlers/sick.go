package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func ApiSickDays(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	staffID := vars["staff_id"]
	year := vars["year"]

	// Validate inputs
	if staffID == "" || year == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Construct start and end of the year
	layout := "2006-01-02"
	startDate := fmt.Sprintf("%s-01-01", year)
	endDate := fmt.Sprintf("%s-12-31", year)

	// Attempt to parse dates to validate
	_, errStart := time.Parse(layout, startDate)
	_, errEnd := time.Parse(layout, endDate)
	if errStart != nil || errEnd != nil {
		http.Error(w, "Invalid year format", http.StatusBadRequest)
		return
	}

	var sickDays []models.Sick
	if err := db.DB.Where("staff_id = ? AND date_from >= ? AND date_from <= ?", staffID, startDate, endDate).
		Find(&sickDays).Error; err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Marshal results to JSON
	jsonResponse, err := json.Marshal(sickDays)
	if err != nil {
		log.Println("JSON Marshalling error:", err)
		http.Error(w, "Failed to process the request", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func ApiSickDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var sickDay struct {
		models.Sick
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	db.DB.Table("sicks").Select("sicks.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on sicks.staff_id = team_members.staff_id").
		Where("sicks.id = ?", param).First(&sickDay)

	json, err := json.Marshal(sickDay)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiSickDash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["staff_id"]

	var sickDash models.SickDash

	const sql = `SELECT
    SUM(CASE WHEN sicks.deducted = 1 THEN sicks.hours ELSE 0 END) AS "sick_days",
    SUM(CASE WHEN sicks.deducted = 0 THEN sicks.hours ELSE 0 END) AS "pending",
	COUNT(*) AS "instances"
FROM
    sicks
WHERE
    sicks.staff_id = ?
GROUP BY
   sicks.staff_id`

	db.DB.Raw(sql, id).Scan(&sickDash)

	json, err := json.Marshal(sickDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
