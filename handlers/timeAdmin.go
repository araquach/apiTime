package handlers

import (
	"encoding/json"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"log"
	"net/http"
)

func ApiTimeAdminDash(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	id := 1

	var timeAdminDash models.TimeAdminDash

	const sql = `SELECT
    (SELECT COUNT(*) FROM holidays h JOIN team_members t ON h.staff_id = t.staff_id WHERE h.approved = 0 AND t.salon = ?) AS holidays,
    (SELECT COUNT(*) FROM free_times f JOIN team_members t ON f.staff_id = t.staff_id WHERE f.approved = 0 AND t.salon = ?) AS free_time,
    (SELECT COUNT(*) FROM lieus l JOIN team_members t ON l.staff_id = t.staff_id WHERE l.approved = 0 AND t.salon = ?) AS lieu_hours,
    (SELECT COUNT(*) FROM sicks s JOIN team_members t ON s.staff_id = s.staff_id WHERE s.deducted= 0 AND t.salon = ?) AS sick_days;`

	db.DB.Raw(sql, id, id, id, id).Scan(&timeAdminDash)

	json, err := json.Marshal(timeAdminDash)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
