package main

import (
	"ApiHealth/controllers"
	"ApiHealth/services"
	"ApiHealth/utils"
	"log"
	"net/http"

	_ "ApiHealth/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API de Monitoreo de Microservicios
// @version 1.0
// @description API para el monitoreo de estado de microservicios.
func main() {
	// Inicializar la conexión a MongoDB
	utils.InitDB()

	// Inicializar monitoreo de servicios registrados
	services.InitMonitoring()

	// Crear el router
	r := mux.NewRouter()

	// Configurar las rutas
	r.HandleFunc("/services", controllers.RegisterService).Methods("POST")
	r.HandleFunc("/services", controllers.GetAllServices).Methods("GET")
	r.HandleFunc("/services/{name}", controllers.GetServiceHealth).Methods("GET")
	r.HandleFunc("/services/{name}", controllers.UpdateService).Methods("PUT")
	r.HandleFunc("/services/{name}", controllers.DeleteService).Methods("DELETE")
	r.HandleFunc("/health", controllers.GetAllHealth).Methods("GET")
	r.HandleFunc("/health/{name}", controllers.GetServiceHealth).Methods("GET")

	r.PathPrefix("/api-docs/").Handler(httpSwagger.WrapHandler)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
