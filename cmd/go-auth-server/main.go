package main

import (
	//"net/http"
	//"os"
	"time"

	"github.com/gavrilaf/go-auth/middleware"
	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	claims := middleware.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authMiddleware := &middleware.AuthMiddleware{Timeout: time.Hour}

	auth := router.Group("/auth")
	{
		auth.POST("/login", authMiddleware.LoginHandler)
	}

	utils := router.Group("/utils")
	utils.Use(authMiddleware.MiddlewareFunc())
	{
		utils.GET("/hello", helloHandler)
	}

	router.Run()
}
