package main

import (
	//"net/http"
	//"os"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/profile"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/gavrilaf/spawn/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	log := logrus.New()

	log.Info("Spawn rest server started")

	router := gin.New()

	router.Use(utils.Logger(log))
	router.Use(gin.Recovery())

	environment := env.GetEnvironment("Test")

	//storage := auth.NewStorageMock(environment)
	apiBridge := api.CreateBridge(environment)
	if apiBridge == nil {
		panic("Can not create storage")
	}

	authMiddleware := auth.CreateMiddleware(apiBridge)

	profileAPI := profile.CreateApi(apiBridge)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authMiddleware.RegisterHandler)
		auth.POST("/login", authMiddleware.LoginHandler)
		auth.POST("/refresh_token", authMiddleware.RefreshHandler)
	}

	profile := router.Group("profile")
	profile.Use(authMiddleware.MiddlewareFunc())
	{
		profile.POST("/confirm_device", profileAPI.ConfirmDevice)
		profile.GET("/", profileAPI.GetUserProfile)

		profile.PUT("/country", profileAPI.UpdateUserCountry)
		profile.PUT("/personal", profileAPI.UpdateUserPersonalInfo)
	}

	router.Run()
}
