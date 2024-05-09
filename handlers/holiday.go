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

func ApiHolidays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["staff_id"]

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var holidays []models.Holiday
	db.DB.Where("staff_id", param).Where("date_from > ? AND date_from < ?", startDate, endDate).Find(&holidays)

	json, err := json.Marshal(holidays)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiHoliday(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var holiday models.Holiday
	db.DB.Where("id", param).Find(&holiday)

	json, err := json.Marshal(holiday)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiHolidayCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var holiday models.Holiday
	var time models.Time

	err := decoder.Decode(&holiday)
	if err != nil {
		panic(err)
	}

	// Start a transaction
	tx := db.DB.Begin()

	res := tx.Where("staff_id", holiday.StaffId).First(&time)
	if res.Error != nil {
		// Handle error, e.g., log it and return
		log.Printf("Error finding booking: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.Create(&holiday)
	if res.Error != nil {
		log.Printf("Error creating holiday: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

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
	vars := mux.Vars(r)
	id := vars["staff_id"]

	var holidayDash models.HolidayDash

	const sql = `SELECT
    times.holiday_ent AS entitlement,
    SUM(CASE WHEN holidays.approved = 0 THEN holidays.requested ELSE 0 END) AS total_pending,
    SUM(CASE WHEN holidays.approved = 1 THEN holidays.requested ELSE 0 END) AS total_booked,
    (times.holiday_ent - SUM(CASE WHEN holidays.approved = 1 THEN holidays.requested ELSE 0 END)) AS remaining,
    (times.holiday_ent - SUM(CASE WHEN holidays.requested != 1 THEN holidays.requested ELSE 0 END)) AS remaining_pending,
    (times.saturday_ent - SUM(CASE WHEN holidays.approved = 0 OR holidays.approved = 1 THEN holidays.saturday ELSE 0 END)) AS sat_pending,
    (times.saturday_ent - SUM(CASE WHEN holidays.approved = 1 THEN holidays.saturday ELSE 0 END)) AS sat_remaining
FROM
    holidays
        JOIN
    times ON holidays.staff_id = times.staff_id
WHERE
    holidays.staff_id = ? AND
    EXTRACT(YEAR FROM holidays.date_from) = EXTRACT(YEAR FROM CURRENT_DATE) AND
    EXTRACT(YEAR FROM holidays.date_to) = EXTRACT(YEAR FROM CURRENT_DATE)
GROUP BY
    times.holiday_ent,
    times.saturday_ent;`

	db.DB.Raw(sql, id).Scan(&holidayDash)

	json, err := json.Marshal(holidayDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
