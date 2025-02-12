package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/chopstickleg/good-code/db"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
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

	_, err = conn.Exec(context.Background(), "INSERT INTO users (email, password) VALUES ($1, $2)", req.Email, req.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User registered"))
}
