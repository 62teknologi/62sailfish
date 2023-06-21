package pubsub_adapter

import (
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(connectionString string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{conn: conn, channel: channel}, nil
}

func (a *RabbitMQ) Publish(_ *gin.Context, topic string, message []byte) error {
	_, err := a.channel.QueueDeclare(topic, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = a.channel.Publish("", topic, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	})
	return err
}

func (a *RabbitMQ) Subscribe(ctx *gin.Context, topic string) (<-chan []byte, error) {
	// Get the context from the request
	c := ctx.Request.Context()

	_, err := a.channel.QueueDeclare(topic, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	msgs, err := a.channel.Consume(topic, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)
	go func() {
		for d := range msgs {
			select {
			case out <- d.Body:
			case <-c.Done():
				close(out)
				return
			}
		}
	}()

	return out, nil
}
