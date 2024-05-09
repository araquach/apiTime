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

func ApiLieuHours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["staff_id"]

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var lieuHours []models.Lieu
	db.DB.Where("staff_id", param).Where("request_date > ? AND request_date < ?", startDate, endDate).Find(&lieuHours)

	json, err := json.Marshal(lieuHours)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiLieuHour(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var lieuHour models.Lieu
	db.DB.Where("id", param).Find(&lieuHour)

	json, err := json.Marshal(lieuHour)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiLieuHourCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var lieu models.Lieu
	var time models.Time

	err := decoder.Decode(&lieu)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()

	res := tx.Where("staff_id", lieu.StaffId).First(&time)
	if res.Error != nil {
		// Handle error, e.g., log it and return
		log.Printf("Error finding time details: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.Create(&lieu)
	if res.Error != nil {
		log.Printf("Error  creating  lieu hour: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

func ApiLieuHourUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var lieu models.Lieu

	err := json.NewDecoder(r.Body).Decode(&lieu)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()

	var originalLieu models.Lieu
	res := tx.First(&originalLieu, id)
	if res.Error != nil {
		log.Printf("Error finding original lieu: %v", res.Error)
		tx.Rollback()
		return
	}

	res = tx.Model(&models.Lieu{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Description": lieu.Description,
		"RequestDate": lieu.RequestDate,
		"Hours":       lieu.Hours,
	})
	if res.Error != nil {
		log.Printf("Error updating lieu: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(lieu)
}

func ApiLieuDash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["staff_id"]

	var lieuDash models.LieuDash

	const sql = `SELECT
    SUM(CASE WHEN lieus.approved = 1 THEN lieus.hours ELSE 0 END) AS "used",
    SUM(CASE WHEN lieus.approved = 0 THEN lieus.hours ELSE 0 END) AS "pending"
FROM
    lieus
WHERE
    lieus.staff_id = ?
GROUP BY
   lieus.staff_id`

	db.DB.Raw(sql, id).Scan(&lieuDash)

	json, err := json.Marshal(lieuDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
