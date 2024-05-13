package handlers

import (
	"encoding/json"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ApiTimeDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["staff_id"]

	var timeInfo models.Time

	db.DB.Where("staff_id", param).First(&timeInfo)

	json, err := json.Marshal(timeInfo)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiTimeDash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["staff_id"]

	var timeDash models.TimeDash

	const sql = `SELECT
    SUM(CASE WHEN holidays.approved = 1 THEN holidays.requested ELSE 0 END) as "holidays",
    SUM(CASE WHEN holidays.approved = 0 THEN holidays.requested ELSE 0 END) as "holidays_pending",
    (SELECT SUM(lieus.hours) FROM lieus WHERE lieus.approved = 1 AND lieus.staff_id = ?) as "lieus",
    (SELECT SUM(lieus.hours) FROM lieus WHERE lieus.approved = 0 AND lieus.staff_id = ?) as "lieu_pending",
    (SELECT SUM(free_times.hours) FROM free_times WHERE free_times.approved = 1 AND free_times.staff_id = ?) as "free_time",
    (SELECT SUM(free_times.hours) FROM free_times WHERE free_times.approved = 0 AND free_times.staff_id = ?) as "free_time_pending",
    (SELECT SUM(sicks.hours) FROM sicks WHERE sicks.deducted = 1 AND sicks.staff_id = ?) as "sick",
    (SELECT SUM(sicks.hours) FROM sicks WHERE sicks.deducted = 0 AND sicks.staff_id = ?) as "sick_pending"
FROM
    holidays
WHERE
    holidays.staff_id = ?`

	db.DB.Raw(sql, id, id, id, id, id, id, id).Scan(&timeDash)

	json, err := json.Marshal(timeDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
