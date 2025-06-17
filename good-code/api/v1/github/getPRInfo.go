package handler

import (
	"encoding/json"
	"net/http"

	db "github.com/chopstickleg/good-code/api/v1/_db"
	middleware "github.com/chopstickleg/good-code/api/v1/_middleware"
	"gorm.io/gorm"
)

func GetPRHandler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods(http.MethodGet)(func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetDB()
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			return
		}

		var req struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		var pr db.AiRoast
		err = conn.Model(&db.AiRoast{}).
			Where("id = ?", req.ID).
			First(&pr).
			Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "PR analysis not found", http.StatusNotFound)
				return
			}
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
	})(w, r)
}
