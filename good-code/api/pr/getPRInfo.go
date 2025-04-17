package handler

import (
	"encoding/json"
	"net/http"

	db "github.com/chopstickleg/good-code/api/_db"
	middleware "github.com/chopstickleg/good-code/api/_middleware"
)

func GetPRHandler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods(http.MethodGet)(func(w http.ResponseWriter, r *http.Request) {
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
		err = conn.Model(&db.Pull_request{}).
			Where("ID = ?", req.ID).
			Find(pr).
			Error
		if err != nil {
			http.Error(w, "Error retrieving data from db", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := json.NewEncoder(w)

		err = response.Encode(pr)
		if err != nil {
			http.Error(w, "Error sending response", http.StatusInternalServerError)
			return
		}
		return
	})(w, r)
}
