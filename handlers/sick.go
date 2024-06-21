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

func ApiSickDays(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["staff_id"]

	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var sickDays []models.Sick
	db.DB.Where("staff_id", param).Where("date_from > ? AND date_from < ?", startDate, endDate).Find(&sickDays)

	json, err := json.Marshal(sickDays)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiSickDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	param := vars["id"]

	var sickDay struct {
		models.Sick
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	db.DB.Table("sicks").Select("sicks.*, team_members.first_name, team_members.last_name").
		Joins("left join team_members on sicks.staff_id = team_members.staff_id").
		Where("sicks.id = ?", param).First(&sickDay)

	json, err := json.Marshal(sickDay)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiSickDash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["staff_id"]

	var sickDash models.SickDash

	const sql = `SELECT
    SUM(CASE WHEN sicks.deducted = 1 THEN sicks.hours ELSE 0 END) AS "sick_days",
    SUM(CASE WHEN sicks.deducted = 0 THEN sicks.hours ELSE 0 END) AS "pending",
	COUNT(*) AS "instances"
FROM
    sicks
WHERE
    sicks.staff_id = ?
GROUP BY
   sicks.staff_id`

	db.DB.Raw(sql, id).Scan(&sickDash)

	json, err := json.Marshal(sickDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
