package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/chopstickleg/good-code/db"
	"github.com/gorilla/mux"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
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

	var hashedPassword string
	err = conn.QueryRow(context.Background(), "SELECT password FROM users WHERE email=$1", req.Email).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// TODO: Compare hashed password and generate JWT token

	w.Write([]byte("Login successful"))
}

func ServeLogin() {
	r := mux.NewRouter()
	r.HandleFunc("/", LoginHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}
