package pubsub_adapter

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type GoogleCloudPubSub struct {
	client *pubsub.Client
}

func NewGoogleCloudPubSub(projectID string) (*GoogleCloudPubSub, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &GoogleCloudPubSub{client: client}, nil
}

func (a *GoogleCloudPubSub) Publish(ctx *gin.Context, topic string, message []byte) error {
	t := a.client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{Data: message})

	_, err := result.Get(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *GoogleCloudPubSub) Subscribe(ctx *gin.Context, topic string) (<-chan []byte, error) {
	messages := make(chan []byte)
	t := a.client.Topic(topic)
	sub := a.client.Subscription(topic)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		_, err := a.client.CreateSubscription(ctx, topic, pubsub.SubscriptionConfig{
			Topic:       t,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}
	go func() {
		err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			messages <- msg.Data
			msg.Ack()
		})
		if err != nil {
			close(messages)
			fmt.Printf("Error receiving message: %v", err)
		}
	}()
	return messages, nil
}
