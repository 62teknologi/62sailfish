package utils

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func PingRabbitMQ(connectionURL string, queueName string, pingMessage string) error {
	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(pingMessage),
		},
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	select {
	case response := <-msgs:
		log.Printf("Received response: %s", response.Body)
	case <-time.After(5 * time.Second):
		return fmt.Errorf("no response received within timeout")
	}

	return nil
}
