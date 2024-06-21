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
	vars := mux.Vars(r)
	id := vars["staff_id"]

	var holidayDash models.HolidayDash

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
        LEFT JOIN holidays ON times.staff_id = holidays.staff_id
WHERE
    times.staff_id = ? AND
    (holidays.date_from IS NULL OR EXTRACT(YEAR FROM holidays.date_from) = EXTRACT(YEAR FROM CURRENT_DATE)) AND
    (holidays.date_to IS NULL OR EXTRACT(YEAR FROM holidays.date_to) = EXTRACT(YEAR FROM CURRENT_DATE))
GROUP BY
    times.holiday_ent,
    times.saturday_ent`

	db.DB.Raw(sql, id).Scan(&holidayDash)

	json, err := json.Marshal(holidayDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
