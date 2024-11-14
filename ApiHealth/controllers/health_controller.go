package controllers

import (
	"ApiHealth/models"
	"ApiHealth/services"
	"ApiHealth/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Registrar un nuevo microservicio (POST /services)
// RegisterService godoc
// @Summary Registra un nuevo servicio
// @Description Registra un nuevo microservicio en el sistema de monitoreo
// @Tags Services
// @Accept  json
// @Produce  json
// @Param service body models.Microservice true "Datos del servicio"
// @Success 201 {object} models.Microservice
// @Router /services [post]
func RegisterService(w http.ResponseWriter, r *http.Request) {
	var service models.Microservice
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Llamar a la función que guarda el servicio en MongoDB y en memoria
	services.RegisterService(service)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(service)
}

// Obtener todos los microservicios registrados (GET /services)
// GetAllServices godoc
// @Summary Obtiene todos los servicios registrados
// @Description Recupera una lista de todos los microservicios monitoreados en el sistema
// @Tags Services
// @Produce  json
// @Success 200 {array} models.Microservice
// @Router /services [get]
func GetAllServices(w http.ResponseWriter, r *http.Request) {
	collection := utils.GetCollection("services")

	// Obtener todos los servicios desde MongoDB
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error al obtener servicios", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var services []models.Microservice
	if err = cursor.All(context.TODO(), &services); err != nil {
		http.Error(w, "Error al parsear resultados", http.StatusInternalServerError)
		return
	}

	// Devolver todos los datos completos de los microservicios
	json.NewEncoder(w).Encode(services)
}

// Obtener datos completos de un microservicio especifico (GET /services/{name})
// GetService godoc
// @Summary Obtiene los datos de un servicio específico
// @Description Recupera los datos completos de un microservicio registrado por su nombre
// @Tags Services
// @Param name path string true "Nombre del servicio"
// @Produce  json
// @Success 200 {object} models.Microservice
// @Failure 404 {string} string "Service not found"
// @Router /services/{name} [get]
func GetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(w, "Nombre del servicio requerido", http.StatusBadRequest)
		return
	}

	collection := utils.GetCollection("services")

	// Buscar el servicio específico en MongoDB
	var service models.Microservice
	filter := bson.M{"name": name}
	err := collection.FindOne(context.TODO(), filter).Decode(&service)
	if err != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(service)
}

// Actualizar datos de un microservicio existente (PUT /services/{name})
// UpdateService godoc
// @Summary Actualiza un servicio existente
// @Description Modifica los datos de un microservicio registrado
// @Tags Services
// @Param name path string true "Nombre del servicio"
// @Param service body models.Microservice true "Datos actualizados del servicio"
// @Produce  json
// @Success 200 {object} models.Microservice
// @Failure 500 {string} string "Error al actualizar el microservicio"
// @Router /services/{name} [put]
func UpdateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var updatedService models.Microservice
	if err := json.NewDecoder(r.Body).Decode(&updatedService); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := services.UpdateService(name, updatedService)
	if err != nil {
		http.Error(w, "Error al actualizar el microservicio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedService)
}

// Eliminar un microservicio (DELETE /services/{name})
// DeleteService godoc
// @Summary Elimina un servicio
// @Description Elimina un microservicio registrado del sistema de monitoreo
// @Tags Services
// @Param name path string true "Nombre del servicio"
// @Produce  json
// @Success 200 {string} string "Microservicio eliminado exitosamente"
// @Failure 500 {string} string "Error al eliminar el microservicio"
// @Router /services/{name} [delete]
func DeleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := services.DeleteService(name)
	if err != nil {
		http.Error(w, "Error al eliminar el microservicio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Microservicio eliminado exitosamente"})
}

// Obtener estado simple de todos los servicios (GET /health)
// GetAllHealth godoc
// @Summary Obtiene el estado de todos los servicios
// @Description Recupera el estado de salud (nombre, URL, y estado) de todos los servicios registrados
// @Tags Health
// @Produce  json
// @Success 200 {array} map[string]string
// @Router /health [get]
func GetAllHealth(w http.ResponseWriter, r *http.Request) {
	collection := utils.GetCollection("services")

	// Obtener todos los servicios desde MongoDB
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error al obtener servicios", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var services []models.Microservice
	if err = cursor.All(context.TODO(), &services); err != nil {
		http.Error(w, "Error al parsear resultados", http.StatusInternalServerError)
		return
	}

	// Devolver solo el nombre, URL y estado
	var healthList []map[string]string
	for _, service := range services {
		healthList = append(healthList, map[string]string{
			"name":   service.Name,
			"url":    service.Endpoint,
			"status": service.Status,
		})
	}

	json.NewEncoder(w).Encode(healthList)
}

// Obtener estado de un servicio específico (GET /health/{name})
// GetServiceHealth godoc
// @Summary Obtiene el estado de un servicio específico
// @Description Recupera el estado de salud de un servicio registrado por su nombre
// @Tags Health
// @Param name path string true "Nombre del servicio"
// @Produce  json
// @Success 200 {object} models.Microservice
// @Failure 404 {string} string "Service not found"
// @Router /health/{name} [get]
func GetServiceHealth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	collection := utils.GetCollection("services")

	// Buscar el servicio específico en MongoDB
	var service models.Microservice
	filter := bson.M{"name": name}
	err := collection.FindOne(context.TODO(), filter).Decode(&service)
	if err != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(service)
}
