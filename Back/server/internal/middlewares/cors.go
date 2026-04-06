package middlewares

import (
	"net/http"

	"server/internal/domain/models"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response to Angular
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
