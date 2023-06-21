package main

import (
	"fmt"

	"github.com/62teknologi/62sailfish/62golib/utils"
	"github.com/62teknologi/62sailfish/app/http/controllers"
	"github.com/62teknologi/62sailfish/app/http/middlewares"
	"github.com/62teknologi/62sailfish/app/interfaces"
	"github.com/62teknologi/62sailfish/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var notificationController controllers.Notification
var firebaseController controllers.Firebase
var chatController controllers.Chat

func main() {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)
		return
	}

	utils.ConnectDatabase(loadedConfig.DBDriver, loadedConfig.DBSource, "")
	utils.InitPluralize()

	r := gin.Default()

	// Setup file static
	r.Static("/public", "./public")

	//allowedIPs := strings.Split(loadedConfig.IpWhitelist, ",")

	// Use the IPWhitelist middleware
	//r.Use(middlewares.IPWhitelist(allowedIPs))

	// Enable CORS for requests from localhost
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("health", notificationController.Health).Use(middlewares.DbSelectorMiddleware())

	apiV1 := r.Group("/api/v1").Use(middlewares.DbSelectorMiddleware())
	{
		RegisterRoute(apiV1, "notifications", &controllers.Notification{})
		FirebaseRoute(apiV1, "token")
		ChatRoute(apiV1, "chats")
	}

	err = r.Run(loadedConfig.HTTPServerAddress)

	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}

func RegisterRoute(r gin.IRoutes, t string, c interfaces.Crud) {
	r.GET("/"+t+"/:id", c.Find)
	r.GET("/"+t+"", c.FindAll)
	r.POST("/"+t+"", c.Create)
	r.POST("/"+t+"/push", notificationController.Push)
	r.GET("/"+t+"/consume", notificationController.Consume)
	r.PUT("/"+t+"/:id", c.Update)
	r.DELETE("/"+t+"/:id", c.Delete)
}

func FirebaseRoute(r gin.IRoutes, t string) {
	r.POST("/"+t+"/push", firebaseController.PostToken)
	r.POST("/"+t+"/push-notifications", firebaseController.PushNotification)
	r.POST("/"+t+"/subscribe-topic", firebaseController.SubscribeTopic)
	r.POST("/"+t+"/unsubscribe-topic", firebaseController.UnsubscribeTopic)
	r.POST("/"+t+"/post-topic", firebaseController.PushTopicNotification)
}

func ChatRoute(r gin.IRoutes, t string) {
	r.POST("/"+t+"/send", chatController.Send)
}
