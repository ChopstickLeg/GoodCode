package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	dataJson := map[string]interface{}{
		"data": map[string]interface{}{
			"message": "Hello World!",
		},
	}
	err := json.NewEncoder(w).Encode(dataJson)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func Serve() {
	r := mux.NewRouter()
	r.HandleFunc("/", Handler).Methods("GET")
	http.ListenAndServe(":8080", r)
}
