package main

import (
	"ApiGateway/config"
	"ApiGateway/internal/handlers"
	"ApiGateway/internal/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	// Conexión a RabbitMQ
	conn, err := utils.ConnectRabbitMQ(cfg.RabbitMQServiceURL)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening a channel: %v", err)
	}
	defer ch.Close()

	// Usar mux como enrutador
	r := mux.NewRouter()

	// Rutas con manejo de parámetros dinámicos
	r.HandleFunc("/auth", handlers.AuthHandler).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/profile/{id_user}", handlers.GetProfileHandler).Methods("GET")
	r.HandleFunc("/update-profile", handlers.UpdateProfileHandler).Methods("PUT")
	r.HandleFunc("/health", handlers.GeneralHealthHandler).Methods("GET")
	r.HandleFunc("/docs", handlers.OpenAPIDocHandler).Methods("GET")

	// Aplicar CORS si es necesario
	corsHandler := handlers.EnableCors(r)

	// Iniciar el servidor en el puerto 5000
	log.Println("Servidor Gateway Iniciado en el puerto 5000")
	log.Fatal(http.ListenAndServe(":5000", corsHandler))
}
