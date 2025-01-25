package handlers

import (
	"encoding/json"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func ApiSickDays(w http.ResponseWriter, r *http.Request) {
	// Extract staff ID from the route variables
	vars := mux.Vars(r)
	staffID := vars["staff_id"]

	// Validate input
	if staffID == "" {
		http.Error(w, "Missing required parameter: staff_id", http.StatusBadRequest)
		return
	}

	// Calculate rolling 12-month date range
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 12 months back

	// Fetch sick days within the rolling 12-month period
	var sickDays []models.Sick
	if err := db.DB.Where("staff_id = ? AND date_from BETWEEN ? AND ?", staffID, startDate, now).
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
	w.Header().Set("Content-Type", "application/json")
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
	// Extract staff ID from the route variables
	vars := mux.Vars(r)
	id := vars["staff_id"]

	// Get the current date and calculate the start date (12 months back)
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 12 months back

	// Define the SQL query, filtering by the rolling 12-month range
	const sql = `SELECT
		SUM(CASE WHEN sicks.deducted = 1 THEN sicks.hours ELSE 0 END) AS sick_days,
		SUM(CASE WHEN sicks.deducted = 0 THEN sicks.hours ELSE 0 END) AS pending,
		COUNT(*) AS instances
	FROM
		sicks
	WHERE
		sicks.staff_id = ? AND
		sicks.date_from BETWEEN ? AND ?
	GROUP BY
		sicks.staff_id`

	// Prepare to store the query result
	var sickDash models.SickDash

	// Execute the query, substituting the staff ID and date range
	if err := db.DB.Raw(sql, id, startDate, now).Scan(&sickDash).Error; err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
		return
	}

	// Marshal the result into JSON
	jsonData, err := json.Marshal(sickDash)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	// Set headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
