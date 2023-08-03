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
	db.DB.Where("staff_id", param).Where("request_date_from > ? AND request_date_from < ?", startDate, endDate).Find(&holidays)

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

	time.SaturdaysPending += holiday.Saturday
	time.HolidaysPending += holiday.HoursRequested

	res = tx.Model(&time).UpdateColumns(map[string]interface{}{"saturdays_pending": time.SaturdaysPending, "holidays_pending": time.HolidaysPending})

	if res.Error != nil {
		// Handle error, e.g., log it and return
		log.Printf("Error updating booking: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.Create(&holiday)
	if res.Error != nil {
		log.Printf("Error creating holiday: %v", res.Error)
		tx.Rollback()
		return
	}

	// If everything went well, commit the transaction
	tx.Commit()

	return
}

func ApiHolidayUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var holiday models.Holiday
	var time models.Time

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

	// Find the time entry
	res = tx.First(&time)
	if res.Error != nil {
		log.Printf("Error finding time: %v", res.Error)
		tx.Rollback()
		return
	}

	// Calculate the differences
	saturdaysDiff := holiday.Saturday - originalHoliday.Saturday
	holidaysDiff := holiday.HoursRequested - originalHoliday.HoursRequested

	// Update the time entry
	time.SaturdaysPending += saturdaysDiff
	time.HolidaysPending += holidaysDiff

	res = tx.Model(&time).UpdateColumns(map[string]interface{}{"saturdays_pending": time.SaturdaysPending, "holidays_pending": time.HolidaysPending})
	if res.Error != nil {
		log.Printf("Error updating time: %v", res.Error)
		tx.Rollback()
		return
	}

	// Update the holiday in the database
	res = tx.Model(&models.Holiday{}).Where("id = ?", id).Updates(holiday)
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
