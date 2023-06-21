package fcm

import (
	"context"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	utils2 "github.com/62teknologi/62sailfish/app/utils"
	"github.com/62teknologi/62sailfish/config"
	"google.golang.org/api/option"
)

type FirebaseNotificationAdapter struct {
	app *firebase.App
}

func NewFirebaseNotificationAdapter() (*FirebaseNotificationAdapter, error) {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		//fmt.Printf("cannot load loadedConfig: %w", err)
		return nil, err
	}
	serviceAccountKeyFilePath, err := filepath.Abs(loadedConfig.ServiceKeyPath)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return &FirebaseNotificationAdapter{app: app}, nil
}

func (f *FirebaseNotificationAdapter) PostNotifications(tokens []string, input map[string]any) ([]string, error) {
	ctx := context.Background()
	client, err := f.app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	data, ec := utils2.ConvertJsonDataAnyToString(input)
	if ec != nil {
		return nil, ec
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    input["title"].(string),
			Body:     input["body"].(string),
			ImageURL: input["image"].(string),
		},
		Data:   data,
		Tokens: tokens,

		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: input["title"].(string),
				Body:  input["body"].(string),
				Image: input["image"].(string),
			},
			Data: data,
		},
	}

	response, err := client.SendMulticast(ctx, message)
	if response.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range response.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, tokens[idx])
			}
		}
		return failedTokens, nil
	}

	return nil, nil
}

func (f *FirebaseNotificationAdapter) PostNotification(token string, input map[string]any) error {
	ctx := context.Background()
	client, err := f.app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	data, err := utils2.ConvertJsonDataAnyToString(input)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    input["title"].(string),
			Body:     input["body"].(string),
			ImageURL: input["image"].(string),
		},
		Data:  data,
		Token: token,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: input["title"].(string),
				Body:  input["body"].(string),
				Image: input["image"].(string),
			},
			Data: data,
		},
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (f *FirebaseNotificationAdapter) SubscribeTopic(tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	ctx := context.Background()
	client, err := f.app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	response, errTopic := client.SubscribeToTopic(ctx, tokens, topic)

	if errTopic != nil {
		return nil, errTopic
	}

	return response, nil
}

func (f *FirebaseNotificationAdapter) UnsubscribeTopic(tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	ctx := context.Background()
	client, err := f.app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	response, errTopic := client.UnsubscribeFromTopic(ctx, tokens, topic)

	if errTopic != nil {
		return nil, errTopic
	}

	return response, nil
}

func (f *FirebaseNotificationAdapter) PostTopic(topic string, input map[string]any) (string, error) {
	ctx := context.Background()
	client, err := f.app.Messaging(ctx)
	if err != nil {
		return "", err
	}

	data, ec := utils2.ConvertJsonDataAnyToString(input)
	if ec != nil {
		return "", ec
	}

	message := &messaging.Message{
		Data:  data,
		Topic: topic,
		Notification: &messaging.Notification{
			Title:    input["title"].(string),
			Body:     input["body"].(string),
			ImageURL: input["image"].(string),
		},
	}

	res, err := client.Send(ctx, message)
	if err != nil {
		return "", err
	}

	return res, nil
}
