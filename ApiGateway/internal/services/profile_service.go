package services

import (
	"ApiGateway/config"
	"ApiGateway/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// GetProfile obtiene los datos del perfil desde el servicio de gestión de perfiles
func GetProfile(idUser, token string) (*models.ProfileData, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.LoadConfig().ProfileServiceURL+"/profile/"+"carlos", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	log.Printf("Sending request to Profile Service for user ID: %s with token: %s", idUser, token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Profile service responded with error: %s", string(body))
		return nil, errors.New("failed to retrieve profile data: " + resp.Status)
	}

	var profileData models.ProfileData
	if err := json.NewDecoder(resp.Body).Decode(&profileData); err != nil {
		return nil, err
	}

	log.Printf("ProfileData retrieved: %+v", profileData)
	return &profileData, nil
}

// UpdateProfile actualiza los datos del perfil en el servicio de gestión de perfiles
func UpdateProfile(profileData models.ProfileData) (*models.ProfileData, error) {
	requestBody, _ := json.Marshal(profileData)
	req, err := http.NewRequest("PUT", config.LoadConfig().ProfileServiceURL+"/profile/"+profileData.ID, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+profileData.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to update profile data: " + resp.Status)
	}

	var updatedProfile models.ProfileData
	if err := json.NewDecoder(resp.Body).Decode(&updatedProfile); err != nil {
		return nil, err
	}

	return &updatedProfile, nil
}
