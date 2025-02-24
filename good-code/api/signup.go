package handler

import (
	"encoding/json"
	"net/http"

	"github.com/chopstickleg/good-code/db"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var count int64
	err = conn.Model(&db.User{}).
		Where("email = ?", req.Email).
		Count(&count).
		Error
	if err != nil {
		http.Error(w, "Error getting existing users from db", http.StatusInternalServerError)
	}

	if count > 0 {
		http.Error(w, "Email already exists", http.StatusBadRequest)
	}

	passByte := []byte(req.Password)
	hashByte, err := bcrypt.GenerateFromPassword(passByte, 1)

	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := db.User{
		Email:    req.Email,
		Password: hashByte,
		Enabled:  true,
	}

	err = conn.Create(user).Error

	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := json.NewEncoder(w)

	err = response.Encode(map[string]bool{"success": true})
	if err != nil {
		http.Error(w, "Failed to send response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
