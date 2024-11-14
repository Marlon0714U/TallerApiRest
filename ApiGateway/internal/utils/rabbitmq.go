package utils

import "github.com/streadway/amqp"

func ConnectRabbitMQ(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func PublishLog(ch *amqp.Channel, queueName string, message string) error {
	return ch.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}
