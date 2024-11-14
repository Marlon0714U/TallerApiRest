package services

import (
	"ApiHealth/models"
	"ApiHealth/utils"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var services = make(map[string]*models.Microservice)
var mu sync.Mutex
var previousStatuses = make(map[string]string)

// Inicializa el monitoreo de todos los servicios registrados en MongoDB
func InitMonitoring() {
	collection := utils.GetCollection("services")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error al cargar servicios desde MongoDB:", err)
		return
	}
	defer cursor.Close(context.TODO())

	var loadedServices []models.Microservice
	if err = cursor.All(context.TODO(), &loadedServices); err != nil {
		log.Println("Error al parsear servicios:", err)
		return
	}

	// Iniciar el monitoreo para cada servicio cargado
	for _, service := range loadedServices {
		log.Printf("Iniciando monitoreo para el servicio: %s", service.Name)
		mu.Lock()
		services[service.Name] = &service
		mu.Unlock()
		go MonitorService(service)
	}
}

// Registrar un microservicio para monitorear y guardarlo en MongoDB
func RegisterService(service models.Microservice) {
	collection := utils.GetCollection("services")

	_, err := collection.InsertOne(context.TODO(), service)
	if err != nil {
		log.Println("Error al registrar el microservicio en MongoDB:", err)
		return
	}

	// También agregar al mapa en memoria para monitoreo
	mu.Lock()
	services[service.Name] = &service
	mu.Unlock()

	go MonitorService(service)
}

// Actualizar un microservicio en MongoDB y en memoria
func UpdateService(name string, updatedService models.Microservice) error {
	collection := utils.GetCollection("services")

	filter := bson.M{"name": name}
	update := bson.M{
		"$set": bson.M{
			"endpoint":  updatedService.Endpoint,
			"frequency": updatedService.Frequency,
			"emails":    updatedService.Emails,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error al actualizar el microservicio en MongoDB:", err)
		return err
	}

	// Actualizar en memoria
	mu.Lock()
	if existingService, ok := services[name]; ok {
		existingService.Endpoint = updatedService.Endpoint
		existingService.Frequency = updatedService.Frequency
		existingService.Emails = updatedService.Emails
	}
	mu.Unlock()

	return nil
}

// Eliminar un microservicio en MongoDB y detener el monitoreo
func DeleteService(name string) error {
	collection := utils.GetCollection("services")

	// Eliminar de MongoDB
	filter := bson.M{"name": name}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println("Error al eliminar el microservicio en MongoDB:", err)
		return err
	}

	// Eliminar de memoria
	mu.Lock()
	delete(services, name)
	mu.Unlock()

	return nil
}

// Monitorear la salud de los microservicios
func MonitorService(service models.Microservice) {
	ticker := time.NewTicker(time.Duration(service.Frequency) * time.Second)
	for range ticker.C {
		log.Printf("Verificando el estado del microservicio: %s en %s", service.Name, time.Now().Format("2006-01-02 15:04:05"))
		CheckHealth(service)
	}
}

// Verificar el estado de un microservicio y enviar notificación si cambia
func CheckHealth(service models.Microservice) {
	resp, err := http.Get(service.Endpoint)
	mu.Lock()
	defer mu.Unlock()

	currentStatus := "UP"
	if err != nil || resp.StatusCode != http.StatusOK {
		currentStatus = "DOWN"
	}

	// Verificar si el estado ha cambiado a DOWN
	if previousStatuses[service.Name] == "UP" && currentStatus == "DOWN" {
		log.Printf("Cambio de estado detectado en %s: %s -> %s", service.Name, previousStatuses[service.Name], currentStatus)

		// Crear mensaje de notificación y enviarlo directamente a la API de notificaciones
		message := models.NotificationMessage{
			Tipo:         "email",
			Mensaje:      service.Name + " está " + currentStatus,
			Destinatario: service.Emails[0],
		}
		log.Printf("Json: %s", message)
		sendNotificationDirectly(message)

		messagesms := models.NotificationMessage{
			Tipo:         "sms",
			Mensaje:      service.Name + " está " + currentStatus,
			Destinatario: "+573162379799",
		}
		sendNotificationDirectly(messagesms)

	}

	// Actualizar el estado en memoria y en MongoDB
	previousStatuses[service.Name] = currentStatus
	services[service.Name].Status = currentStatus
	UpdateServiceStatus(service.Name, currentStatus)
}

// Actualizar el estado de un microservicio en MongoDB
func UpdateServiceStatus(name, status string) {
	collection := utils.GetCollection("services")
	filter := bson.M{"name": name}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error al actualizar el estado en MongoDB:", err)
	}
}

// Función para enviar notificación directamente a la API de notificaciones
func sendNotificationDirectly(message models.NotificationMessage) {
	notificationsAPIURL := os.Getenv("NOTIFICATIONS_API_URL")
	endpoint := notificationsAPIURL + "/send"

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje de notificación: %s", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al crear la solicitud de notificación: %s", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error al enviar notificación a la API de notificaciones: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Error: La API de notificaciones respondió con el código de estado %d. Detalle: %s", resp.StatusCode, string(bodyBytes))
	} else {
		log.Println("Notificación enviada correctamente a la API de notificaciones.")
	}
}
