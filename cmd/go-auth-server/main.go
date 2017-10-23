package main

import (
	//"net/http"
	//"os"
	"github.com/gavrilaf/go-auth/middleware"
	"github.com/gavrilaf/go-auth/storage"
	"github.com/gin-gonic/gin"
	"time"
)

func helloHandler(c *gin.Context) {
	userId := c.GetString(middleware.UserIDName)
	clientId := c.GetString(middleware.ClientIDName)

	c.JSON(200, gin.H{
		"user_id":   userId,
		"client_id": clientId,
	})
}

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	clientsStorage := storage.ClientsStorageMock{}
	usersStorage := storage.UsersStorageMock{}
	sessionsStorage := storage.NewMemorySessionsStorage()

	storage := storage.StorageFacade{Clients: clientsStorage, Users: usersStorage, Sessions: sessionsStorage}
	authMiddleware := &middleware.AuthMiddleware{Timeout: time.Minute, MaxRefresh: time.Hour, Storage: storage}

	auth := router.Group("/auth")
	{
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
