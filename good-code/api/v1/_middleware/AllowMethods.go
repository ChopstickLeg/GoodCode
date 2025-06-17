package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func AllowMethods(allowed ...string) func(http.HandlerFunc) http.HandlerFunc {
	if len(allowed) == 0 {
		log.Printf("Warning: AllowMethods called with no allowed methods")
	}

	allowedMap := make(map[string]bool)
	for _, m := range allowed {
		if m == "" {
			log.Printf("Warning: Empty method string provided to AllowMethods")
			continue
		}
		allowedMap[m] = true
	}

	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !allowedMap[r.Method] {
				log.Printf("Method %s not allowed for %s, allowed methods: %v", r.Method, r.URL.Path, allowed)
				http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
				return
			}
			log.Printf("Allowed method %s for %s", r.Method, r.URL.Path)
			hf(w, r)
		}
	}
}
