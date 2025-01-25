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

func ApiLieuHours(w http.ResponseWriter, r *http.Request) {
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

	// Fetch lieu hours within the rolling 6-month period
	var lieuHours []models.Lieu
	if err := db.DB.Where("staff_id = ? AND request_date BETWEEN ? AND ?", staffID, startDate, now).
		Find(&lieuHours).Error; err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Marshal results to JSON
	jsonResponse, err := json.Marshal(lieuHours)
	if err != nil {
		log.Println("JSON Marshalling error:", err)
		http.Error(w, "Failed to process the request", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func ApiLieuHour(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var lieuHourDetail struct {
		models.Lieu
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	db.DB.Table("lieus").Select("lieus.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on lieus.staff_id = team_members.staff_id").
		Where("lieus.id = ?", param).First(&lieuHourDetail)

	json, err := json.Marshal(lieuHourDetail)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiLieuHourCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var lieu models.Lieu

	err := decoder.Decode(&lieu)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := db.DB.Create(&lieu)
	if result.Error != nil {
		http.Error(w, "Failed to create record: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}
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
