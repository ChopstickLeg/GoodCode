package handler

import (
	"encoding/json"
	"net/http"

	db "github.com/chopstickleg/good-code/api/v1/_db"
	middleware "github.com/chopstickleg/good-code/api/v1/_middleware"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods("POST")(func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetDB()
		if err != nil {
			http.Error(w, "Failed to connect to database: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var req struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		var user db.User_login
		err = conn.Model(&db.User_login{}).
			Where("email = ?", req.Email).
			Find(&user).
			Error
		if err != nil {
			http.Error(w, "Error querying DB: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if user.Enabled {
			http.Error(w, "User already exists", http.StatusUnauthorized)
			return
		}

		passByte := []byte(req.Password)
		hashByte, err := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)

		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		user = db.User_login{
			Email:    req.Email,
			Name:     req.Name,
			Password: hashByte,
			Enabled:  true,
		}

		err = conn.Create(&user).Error

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
	})(w, r)
}
