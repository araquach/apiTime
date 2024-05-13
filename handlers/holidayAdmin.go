package handlers

import (
	"encoding/json"
	"github.com/araquach/apiTime/models"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ApiHolidayApprove(w http.ResponseWriter, r *http.Request) {
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
		"Approved": holiday.Approved,
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
