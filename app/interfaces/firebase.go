package interfaces

import "firebase.google.com/go/messaging"

type FirebasePushNotification interface {
	//PostToken(ctx *gin.Context)
	PostNotification(token string, input map[string]any) error
	PostNotifications(tokens []string, input map[string]any) ([]string, error)
	SubscribeTopic(tokens []string, topic string) (*messaging.TopicManagementResponse, error)
	UnsubscribeTopic(tokens []string, topic string) (*messaging.TopicManagementResponse, error)
	PostTopic(topic string, input map[string]any) (string, error)
}
