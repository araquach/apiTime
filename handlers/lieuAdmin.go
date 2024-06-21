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

func ApiLieuAdminDash(w http.ResponseWriter, r *http.Request) {

	var lieuAdminDash models.LieuAdminDash

	const sql = `SELECT COUNT(*) AS pending FROM lieus WHERE approved = 0;`

	db.DB.Raw(sql).Scan(&lieuAdminDash)

	json, err := json.Marshal(lieuAdminDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiAdminLieuPending(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Result struct {
		models.Lieu
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var results []Result

	db.DB.Table("lieus").Select("lieus.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on team_members.staff_id = lieus.staff_id").
		Where("lieus.request_date > ? AND lieus.request_date < ?", startDate, endDate).
		Where("lieus.approved = ?", 0).
		Scan(&results)

	json, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiAdminLieuHours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["staff_id"]

	type Result struct {
		models.Lieu
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var results []Result

	db.DB.Table("lieus").Where("lieus.staff_id", id).Select("lieus.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on team_members.staff_id = lieus.staff_id").
		Where("lieus.request_date > ? AND lieus.request_date < ?", startDate, endDate).
		Scan(&results)

	json, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiLieuApprove(w http.ResponseWriter, r *http.Request) {
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
		"Approved": lieu.Approved,
	})
	if res.Error != nil {
		log.Printf("Error updating lieu: %v", res.Error)
		tx.Rollback()
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(lieu)
}
