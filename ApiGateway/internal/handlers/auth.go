package handlers

import (
	"ApiGateway/internal/models"
	"ApiGateway/internal/services"
	"encoding/json"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var request models.RegisterData
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	authData, err := services.AuthUser(request.Username, request.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authData)
}

func EnableCors(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {

			w.WriteHeader(http.StatusOK)

			return

		}

		next.ServeHTTP(w, r)

	})

}
