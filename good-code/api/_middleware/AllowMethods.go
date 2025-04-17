package middleware

import (
	"fmt"
	"net/http"
)

func AllowMethods(allowed ...string) func(http.HandlerFunc) http.HandlerFunc {
	allowedMap := make(map[string]bool)
	for _, m := range allowed {
		allowedMap[m] = true
	}

	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !allowedMap[r.Method] {
				http.Error(w, fmt.Sprint("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
				return
			}
			hf(w, r)
		}
	}
}
