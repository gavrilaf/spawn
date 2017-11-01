package main

import (
	//"net/http"
	//"os"
	"time"

	"github.com/gavrilaf/spawn/pkg/auth"
	"github.com/gavrilaf/spawn/pkg/ginlog"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func helloHandler(c *gin.Context) {
	userId := c.GetString(auth.UserIDName)
	clientId := c.GetString(auth.ClientIDName)

	c.JSON(200, gin.H{
		"user_id":   userId,
		"client_id": clientId,
	})
}

func main() {

	log := logrus.New()

	log.Info("Spawn rest server started")

	router := gin.New()

	router.Use(ginlog.Logger(log))
	router.Use(gin.Recovery())

	storage := auth.StorageFacade{Clients: auth.NewClientsStorageMock(), Users: auth.NewUsersStorageMock(), Sessions: auth.NewMemorySessionsStorage()}
	authMiddleware := &auth.Middleware{Timeout: time.Minute, MaxRefresh: time.Hour, Storage: storage, Log: log}

	auth := router.Group("/auth")
	{
		auth.POST("/register", authMiddleware.RegisterHandler)
		auth.POST("/login", authMiddleware.LoginHandler)
		auth.POST("/refresh_token", authMiddleware.RefreshHandler)
	}

	utils := router.Group("/utils")
	utils.Use(authMiddleware.MiddlewareFunc())
	{
		utils.GET("/hello", helloHandler)
	}

	router.Run()
}
