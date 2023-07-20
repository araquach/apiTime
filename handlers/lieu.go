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
	db.DB.Where("staff_id", param).Where("date_regarding > ? AND date_regarding < ?", startDate, endDate).Find(&lieuHours)

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
		panic(err)
	}

	tx := db.DB.Begin()

	res := tx.First(&time)
	if res.Error != nil {
		// Handle error, e.g., log it and return
		log.Printf("Error finding time details: %v", res.Error)
		tx.Rollback()
		return
	}

	time.LieuHours += lieu.LieuHours

	res = tx.Model(&time).UpdateColumn("lieu_hours", time.LieuHours)
	if res.Error != nil {
		// Handle error, e.g., log it and return
		log.Printf("Error updating lieu hours in times: %v", res.Error)
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
