package handler

import (
	"encoding/json"
	"net/http"
)

func verifyJWTHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := json.NewEncoder(w)
	err := response.Encode(map[string]bool{"loggedIn": false})
	if err != nil {
		http.Error(w, "Failed to send response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
