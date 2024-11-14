package utils

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// InitDB inicializa la conexión a MongoDB
func InitDB() {

	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		log.Fatal("La variable de entorno MONGO_URL no está configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verificar la conexión
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Conectado a MongoDB")
	Client = client
}

// GetCollection devuelve una referencia a una colección específica
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("microservice_monitor").Collection(collectionName)
}
