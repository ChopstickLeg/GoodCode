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

	passByte := []byte(req.Password)
	hashByte, err := bcrypt.GenerateFromPassword(passByte, 1)

	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = conn.Exec("INSERT INTO user_login (email, password, enabled) VALUES ($1, $2, TRUE)", req.Email, hashByte)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User registered"))
}
