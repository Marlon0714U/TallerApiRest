package config

import (
    "os"
)

type Config struct {
    AuthServiceURL     string
    ProfileServiceURL  string
    RabbitMQServiceURL string
    QueueName          string
}

func LoadConfig() *Config {
    return &Config{
        AuthServiceURL:     os.Getenv("AUTH_SERVICE_URL"),
        ProfileServiceURL:  os.Getenv("PROFILE_SERVICE_URL"),
        RabbitMQServiceURL: os.Getenv("RABBITMQ_SERVICE_URL"),
        QueueName:          os.Getenv("QUEUE_NAME"),
    }
}
