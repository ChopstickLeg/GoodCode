package handler

import (
	"net/http"
	"os"

	db "github.com/chopstickleg/good-code/api/_db"
)

func MigrateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secret := os.Getenv("MIGRATION_SECRET")
	if secret == "" {
		http.Error(w, "Bad secret", http.StatusInternalServerError)
		return
	}

	passedSecret := r.Header.Get("X-Migration-Secret")
	if passedSecret != secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	err = conn.AutoMigrate(
		&db.UserLogin{},
		&db.Repository{},
		&db.UserRepositoryCollaborator{},
		&db.AiRoast{},
	)
	if err != nil {
		http.Error(w, "Migration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Migration completed successfully"))
}
