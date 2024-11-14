package models

type Microservice struct {
	Name      string   `json:"name"`
	Endpoint  string   `json:"endpoint"`
	Frequency int      `json:"frequency"` // Frecuencia en segundos
	Emails    []string `json:"emails"`
	Status    string   `json:"status"` // "UP" o "DOWN"
}
