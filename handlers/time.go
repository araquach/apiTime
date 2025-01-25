package handlers

import (
	"encoding/json"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func ApiTimeDetails(w http.ResponseWriter, r *http.Request) {
	// Get the staff_id from the request
	vars := mux.Vars(r)
	param := vars["staff_id"]

	var timeInfo models.Time

	// Query the Time record and preload the most recent Schedule
	err := db.DB.Where("staff_id = ?", param).
		Preload("Schedules", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		First(&timeInfo).Error

	if err != nil {
		// Handle error if Time record is not found
		if gorm.ErrRecordNotFound == err {
			http.Error(w, "Time record not found", http.StatusNotFound)
			return
		}
		log.Println("Error querying time record:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal the result to JSON
	jsonData, err := json.Marshal(timeInfo)
	if err != nil {
		log.Println("Error marshalling time info:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func ApiTimeDash(w http.ResponseWriter, r *http.Request) {
	// Extract staff ID and year from the route variables
	vars := mux.Vars(r)
	id := vars["staff_id"]
	year := vars["year"]

	// Calculate the rolling 12-month date range for sick time
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 12 months back

	// Define the SQL query with updated filters
	const sql = `
    SELECT
        SUM(CASE WHEN holidays.approved = 1 THEN holidays.requested ELSE 0 END) as "holidays",
        SUM(CASE WHEN holidays.approved = 0 THEN holidays.requested ELSE 0 END) as "holidays_pending",
        (SELECT SUM(lieus.hours)
         FROM lieus 
         WHERE lieus.approved = 1 AND lieus.staff_id = ?) as "lieus",
        (SELECT SUM(lieus.hours)
         FROM lieus 
         WHERE lieus.approved = 0 AND lieus.staff_id = ?) as "lieu_pending",
        (SELECT SUM(free_times.hours) 
         FROM free_times 
         WHERE free_times.approved = 1 AND free_times.staff_id = ?) as "free_time",
        (SELECT SUM(free_times.hours) 
         FROM free_times 
         WHERE free_times.approved = 0 AND free_times.staff_id = ?) as "free_time_pending",
        (SELECT SUM(sicks.hours) 
         FROM sicks 
         WHERE sicks.deducted = 1 AND sicks.staff_id = ? AND sicks.date_from BETWEEN ? AND ?) as "sick",
        (SELECT SUM(sicks.hours) 
         FROM sicks 
         WHERE sicks.deducted = 0 AND sicks.staff_id = ? AND sicks.date_from BETWEEN ? AND ?) as "sick_pending"
    FROM
        holidays
    WHERE
        holidays.staff_id = ? AND EXTRACT(YEAR FROM holidays.created_at) = ?`

	// Prepare to store the query result
	var timeDash models.TimeDash

	// Execute the query with the updated parameters
	if err := db.DB.Raw(sql, id, id, id, id, id, startDate, now, id, startDate, now, id, year).Scan(&timeDash).Error; err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
		return
	}

	// Marshal the result into JSON
	jsonData, err := json.Marshal(timeDash)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	// Set headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
