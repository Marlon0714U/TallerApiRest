package services

import (
	"ApiGateway/config"
	"ApiGateway/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// AuthUser authenticates the user and returns the token wrapped in a JSON structure.
func AuthUser(email, password string) (map[string]string, error) {
	requestBody, _ := json.Marshal(map[string]string{
		"username": email,
		"password": password,
	})

	authServiceURL := config.LoadConfig().AuthServiceURL + "/login"
	log.Printf("Sending request to auth service at %s", authServiceURL)

	resp, err := http.Post(authServiceURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error connecting to auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Received non-OK response: %s\nResponse body: %s", resp.Status, string(bodyBytes))
		return nil, errors.New("authentication failed")
	}

	var authData models.AuthData
	if err := json.NewDecoder(resp.Body).Decode(&authData); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Return the token wrapped in a JSON structure
	return map[string]string{"token": authData.Token}, nil
}

// GetAuthData obtiene los datos de autenticación del usuario desde el servicio de autenticación
func GetAuthData(idUser string, token string) (*models.AuthData, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.LoadConfig().AuthServiceURL+"/users/"+idUser, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve auth data: " + resp.Status)
	}

	var authData models.AuthData
	if err := json.NewDecoder(resp.Body).Decode(&authData); err != nil {
		return nil, err
	}

	log.Printf("AuthData retrieved: %+v", authData)

	return &authData, nil
}

// UpdateAuthData actualiza los datos de autenticación en el servicio de autenticación
func UpdateAuthData(authData models.AuthData) error {
	requestBody, _ := json.Marshal(authData)
	req, err := http.NewRequest("PUT", config.LoadConfig().AuthServiceURL+"/users/"+strconv.Itoa(authData.IDUser), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authData.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to update auth data: " + resp.Status)
	}
	return nil
}
