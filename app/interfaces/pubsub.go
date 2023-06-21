package interfaces

import (
	"github.com/gin-gonic/gin"
)

type PubSub interface {
	Publish(ctx *gin.Context, topic string, message []byte) error
	Subscribe(ctx *gin.Context, topic string) (<-chan []byte, error)
}
