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

func ApiHolidays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract params from URL
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

	// Query the database
	var holidays []models.Holiday
	if err := db.DB.Where("staff_id = ? AND date_from >= ? AND date_from <= ?", staffID, startDate, endDate).
		Find(&holidays).Error; err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Marshal results to JSON
	jsonResponse, err := json.Marshal(holidays)
	if err != nil {
		log.Println("JSON Marshalling error:", err)
		http.Error(w, "Failed to process the request", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func ApiHoliday(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var holidayDetail struct {
		models.Holiday
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	db.DB.Table("holidays").Select("holidays.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on holidays.staff_id = team_members.staff_id").
		Where("holidays.id = ?", param).First(&holidayDetail)

	json, err := json.Marshal(holidayDetail)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiHolidayCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var holiday models.Holiday

	err := decoder.Decode(&holiday)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := db.DB.Create(&holiday)
	if result.Error != nil {
		http.Error(w, "Failed to create record: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func ApiHolidayUpdate(w http.ResponseWriter, r *http.Request) {
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
		"Requested":   holiday.Requested,
		"Description": holiday.Description,
		"DateFrom":    holiday.DateFrom,
		"DateTo":      holiday.DateTo,
		"Saturday":    holiday.Saturday,
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

func ApiHolidayDash(w http.ResponseWriter, r *http.Request) {
	// Extract staff ID and year from the route variables
	vars := mux.Vars(r)
	id := vars["staff_id"]
	year := vars["year"]

	// Define the SQL query, filtering by the designated year
	const sql = `SELECT
		times.holiday_ent AS entitlement,
		COALESCE(SUM(CASE WHEN holidays.approved = 0 THEN holidays.requested ELSE 0 END), 0) AS total_pending,
		COALESCE(SUM(CASE WHEN holidays.approved = 1 THEN holidays.requested ELSE 0 END), 0) AS total_booked,
		COALESCE((times.holiday_ent - SUM(CASE WHEN holidays.approved = 1 THEN holidays.requested ELSE 0 END)), times.holiday_ent) AS remaining,
		COALESCE((times.holiday_ent - SUM(CASE WHEN holidays.requested != 1 THEN holidays.requested ELSE 0 END)), times.holiday_ent) AS remaining_pending,
		COALESCE((times.saturday_ent - SUM(CASE WHEN holidays.approved = 0 OR holidays.approved = 1 THEN holidays.saturday ELSE 0 END)), times.saturday_ent) AS sat_pending,
		COALESCE((times.saturday_ent - SUM(CASE WHEN holidays.approved = 1 THEN holidays.saturday ELSE 0 END)), times.saturday_ent) AS sat_remaining
	FROM
		times
		LEFT JOIN holidays
		ON times.staff_id = holidays.staff_id
	WHERE
		times.staff_id = ? AND
		(holidays.date_from IS NULL OR EXTRACT(YEAR FROM holidays.date_from) = ?) AND
		(holidays.date_to IS NULL OR EXTRACT(YEAR FROM holidays.date_to) = ?)
	GROUP BY
		times.holiday_ent,
		times.saturday_ent`

	// Prepare to store the query result
	var holidayDash models.HolidayDash

	// Execute the query, substituting the staff ID and year
	if err := db.DB.Raw(sql, id, year, year).Scan(&holidayDash).Error; err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
		return
	}

	// Marshal the result into JSON
	jsonData, err := json.Marshal(holidayDash)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	// Set headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
