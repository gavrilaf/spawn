package main

import (
	//"net/http"
	//"os"
	"time"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/profile"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/gavrilaf/spawn/pkg/ginlog"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	log := logrus.New()

	log.Info("Spawn rest server started")

	router := gin.New()

	router.Use(ginlog.Logger(log))
	router.Use(gin.Recovery())

	environment := env.GetEnvironment("Test")

	//storage := auth.NewStorageMock(environment)
	storage := api.NewBridge(environment)
	if storage == nil {
		panic("Can not create storage")
	}

	authMiddleware := &auth.Middleware{
		Timeout:    time.Minute,
		MaxRefresh: time.Hour,
		Stg:        auth.StorageImpl{StorageBridge: storage},
		Log:        log}

	profileAPI := profile.ProfileApiImpl{StorageBridge: storage}

	auth := router.Group("/auth")
	{
		auth.POST("/register", authMiddleware.RegisterHandler)
		auth.POST("/login", authMiddleware.LoginHandler)
		auth.POST("/refresh_token", authMiddleware.RefreshHandler)
	}

	profile := router.Group("profile")
	profile.Use(authMiddleware.MiddlewareFunc())
	{
		profile.GET("/whoami", profileAPI.WhoAmI)
		profile.POST("/confirm_device", profileAPI.ConfirmDevice)
	}

	//service := router.Group("/service")
	//utils.Use(authMiddleware.MiddlewareFunc())
	//{
	//utils.GET("/whoami", api.WhoAmI)
	//}

	router.Run()
}
