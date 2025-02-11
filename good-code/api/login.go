package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var dbUrl = os.Getenv("TURSO_DB_URL")
var dbAuth = os.Getenv("TURSO_AUTH_TOKEN")

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	conn, err := libsql.Open(dbUrl, libsql.WithAuthToken(dbAuth))
	if err != nil {
		http.Error(w, "DB connection failed", http.StatusInternalServerError)
	}
	defer conn.Close()

	dataJson := map[string]interface{}{
		"data": map[string]interface{}{
			"message": "Hello World!",
		},
	}
	err = json.NewEncoder(w).Encode(dataJson)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
