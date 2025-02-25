package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/chopstickleg/good-code/db"
	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
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

	var user db.User_login
	err = conn.Model(&db.User_login{}).
		Where("email = ?", req.Email).
		Find(&user).
		Error
	if err != nil {
		http.Error(w, "Error querying DB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if user.Enabled == false {
		http.Error(w, "User does not exist", http.StatusUnauthorized)
	}

	incoming := []byte(req.Password)

	matchErr := bcrypt.CompareHashAndPassword(user.Password, incoming)
	if matchErr != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		http.Error(w, "JWT_SECRET_KEY environment variable not set", http.StatusInternalServerError)
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = req.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, "Failed to generate JWT token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := json.NewEncoder(w)

	err = response.Encode(map[string]string{"token": signedToken})
	if err != nil {
		http.Error(w, "Failed to send response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
