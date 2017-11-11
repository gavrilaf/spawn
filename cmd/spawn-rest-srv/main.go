package main

import (
	//"net/http"
	//"os"
	"time"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/auth"
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

	sessionsStorage := auth.NewCacheSessionsStorage(environment)
	if sessionsStorage == nil {
		panic("Can not create sessions cache")
	}

	//usersStorage := auth.NewUsersStorageMock()
	usersStorage := auth.CreateUsersBridge(environment)
	if usersStorage == nil {
		panic("Can not connect to backend")
	}

	storage := auth.StorageFacade{Clients: auth.NewClientsStorageMock(), Users: usersStorage, Sessions: sessionsStorage}
	authMiddleware := &auth.Middleware{Timeout: time.Minute, MaxRefresh: time.Hour, Storage: storage, Log: log}

	api := &api.SpawnApi{Log: log}

	auth := router.Group("/auth")
	{
		auth.POST("/register", authMiddleware.RegisterHandler)
		auth.POST("/login", authMiddleware.LoginHandler)
		auth.POST("/refresh_token", authMiddleware.RefreshHandler)
	}

	utils := router.Group("/service")
	utils.Use(authMiddleware.MiddlewareFunc())
	{
		utils.GET("/whoami", api.WhoAmI)
	}

	router.Run()
}
