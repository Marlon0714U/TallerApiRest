package handlers

import (
	"ApiGateway/internal/models"
	"ApiGateway/internal/services"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// GetProfileHandler maneja las solicitudes para obtener datos completos del usuario
// Handler para obtener el perfil del usuario
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el id_user desde la URL
	vars := mux.Vars(r)
	idUser := vars["id_user"]

	// Obtener el token desde el encabezado de autorización
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	// Quitar "Bearer " del token para usar solo el valor
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Llamada a los servicios de autenticación y perfil
	if _, authErr := services.GetAuthData(idUser, token); authErr != nil {
		http.Error(w, "Failed to retrieve authentication data", http.StatusUnauthorized)
		return
	}

	profileData, profileErr := services.GetProfile(idUser, token)
	if profileErr != nil {
		http.Error(w, "Failed to retrieve profile data", http.StatusInternalServerError)
		return
	}

	// Unificar ambas respuestas
	combinedResponse := map[string]interface{}{
		"id":             idUser,
		"personal_url":   profileData.URLPagina,
		"contact_public": profileData.ContactPublic,
		"address":        profileData.DireccionCorrespondencia,
		"biography":      profileData.Biografia,
		"organization":   profileData.Organizacion,
		"country":        profileData.Pais,
		"nickname":       profileData.Apodo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(combinedResponse)
}

// UpdateProfileHandler maneja las solicitudes para actualizar datos completos del usuario
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var updateData models.CombinedUserUpdateData
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Actualizar datos de autenticación
	authErr := services.UpdateAuthData(updateData.AuthData)
	if authErr != nil {
		http.Error(w, "Failed to update authentication data", http.StatusInternalServerError)
		return
	}

	// Actualizar datos de perfil
	updatedProfile, profileErr := services.UpdateProfile(updateData.ProfileData)
	if profileErr != nil {
		http.Error(w, "Failed to update profile data", http.StatusInternalServerError)
		return
	}

	profileData, profileErr := services.UpdateProfile(updateData.ProfileData)
	if profileErr != nil {
		http.Error(w, "Failed to retrieve profile data", http.StatusInternalServerError)
		return
	}

	// Unificar ambas respuestas
	combinedResponse := map[string]interface{}{
		"id":             updatedProfile.ID,
		"personal_url":   profileData.URLPagina,
		"contact_public": profileData.ContactPublic,
		"address":        profileData.DireccionCorrespondencia,
		"biography":      profileData.Biografia,
		"organization":   profileData.Organizacion,
		"country":        profileData.Pais,
		"nickname":       profileData.Apodo,
	}

	// Responder con éxito
	// Usa updatedProfile en la respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(combinedResponse)
}

// RegisterHandler maneja el registro de usuarios y redirecciona la solicitud al servicio de autenticación
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	usersAPIURL := os.Getenv("AUTH_SERVICE_URL")
	if usersAPIURL == "" {
		log.Println("Variable de entorno USERS_API_URL no definida")
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error al leer el cuerpo de la solicitud:", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	registerEndpoint := usersAPIURL + "/users"
	resp, err := http.Post(registerEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error al realizar la solicitud al servicio de registro:", err)
		http.Error(w, "Error al realizar la solicitud al servicio de registro", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error al leer la respuesta del servicio de registro:", err)
		http.Error(w, "Error al leer la respuesta del servicio de registro", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
