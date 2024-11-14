package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

var (
	startTime      = time.Now().UTC() // Almacenamos el tiempo de inicio del servicio
	serviceVersion = "1.0.0"          // Definimos la versión del servicio
)

type Check struct {
	Data   map[string]string `json:"data"`
	Name   string            `json:"name"`
	Status string            `json:"status"`
}

type HealthResponse struct {
	Status  string  `json:"status"`
	Checks  []Check `json:"checks"`
	Details struct {
		Version string `json:"version"`
		Uptime  string `json:"uptime"`
	} `json:"details"`
}

// Liveness check - /health/live
func LivenessCheckHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).String()
	checkData := map[string]string{
		"from":   time.Now().UTC().Format(time.RFC3339),
		"status": "ALIVE",
	}

	response := HealthResponse{
		Status: "UP",
		Checks: []Check{
			{
				Data:   checkData,
				Name:   "Liveness check",
				Status: "UP",
			},
		},
	}

	response.Details.Version = serviceVersion
	response.Details.Uptime = uptime

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Readiness check - /health/ready
func ReadinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	isDBConnected := true // Simulamos la conexión a la base de datos
	status := "UP"
	if !isDBConnected {
		status = "DOWN"
	}

	uptime := time.Since(startTime).String()
	checkData := map[string]string{
		"from":   time.Now().UTC().Format(time.RFC3339),
		"status": "READY",
	}

	response := HealthResponse{
		Status: status,
		Checks: []Check{
			{
				Data:   checkData,
				Name:   "Readiness check",
				Status: status,
			},
		},
	}

	response.Details.Version = serviceVersion
	response.Details.Uptime = uptime

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// General health status - /health
func GeneralHealthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).String()
	currentTime := time.Now().UTC().Format(time.RFC3339)

	response := HealthResponse{
		Status: "UP",
		Checks: []Check{
			{
				Data:   map[string]string{"from": currentTime, "status": "READY"},
				Name:   "Readiness check",
				Status: "UP",
			},
			{
				Data:   map[string]string{"from": currentTime, "status": "ALIVE"},
				Name:   "Liveness check",
				Status: "UP",
			},
		},
	}

	response.Details.Version = serviceVersion
	response.Details.Uptime = uptime

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
