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

func ApiFreeTimes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract staff ID from the route variables
	vars := mux.Vars(r)
	staffID := vars["staff_id"]

	// Validate input
	if staffID == "" {
		http.Error(w, "Missing required parameter: staff_id", http.StatusBadRequest)
		return
	}

	// Calculate rolling 6-month date range
	now := time.Now()
	startDate := now.AddDate(0, -6, 0) // 6 months back

	// Fetch free times within the rolling 6-month period
	var freeTimes []models.FreeTime
	if err := db.DB.Where("staff_id = ? AND request_date BETWEEN ? AND ?", staffID, startDate, now).
		Find(&freeTimes).Error; err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Marshal results to JSON
	jsonResponse, err := json.Marshal(freeTimes)
	if err != nil {
		log.Println("JSON Marshalling error:", err)
		http.Error(w, "Failed to process the request", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func ApiFreeTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var freeTimeDetail struct {
		models.FreeTime
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	db.DB.Table("free_times").Select("free_times.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on free_times.staff_id = team_members.staff_id").
		Where("free_times.id = ?", param).First(&freeTimeDetail)

	json, err := json.Marshal(freeTimeDetail)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiFreeTimeCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var freeTime models.FreeTime

	err := decoder.Decode(&freeTime)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := db.DB.Create(&freeTime)
	if result.Error != nil {
		http.Error(w, "Failed to create record: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}
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
