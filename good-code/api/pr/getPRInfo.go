package handler

import (
	"encoding/json"
	"net/http"

	db "github.com/chopstickleg/good-code/api/_db"
)

func getPRHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
	}
	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var pr db.Pull_request
	err = conn.Model()
	return
}
