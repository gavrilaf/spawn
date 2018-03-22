package main

import (
	"os"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/account"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/profile"
	"github.com/gavrilaf/spawn/pkg/api/user"
	"github.com/gavrilaf/spawn/pkg/senv"
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

	log.Info("System environment:")
	for _, e := range os.Environ() {
		log.Info(e)
	}

	env := senv.GetEnvironment()
	if env == nil {
		log.Fatal("Could not read environment")
	}

	log.Infof("Web service environment: %s", env.String())

	//storage := auth.NewStorageMock(environment)
	apiBridge := api.CreateBridge(env)
	if apiBridge == nil {
		log.Info("Could not connect to the api bridge")
	}

	authMiddleware := auth.CreateMiddleware(apiBridge)

	profileAPI := profile.CreateApi(apiBridge)
	userAPI := user.CreateApi(apiBridge)
	accountsApi := account.CreateApi(apiBridge)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authMiddleware.RegisterHandler)
		auth.POST("/login", authMiddleware.LoginHandler)
		auth.POST("/refresh_token", authMiddleware.RefreshHandler)
	}

	user := router.Group("user")
	user.Use(authMiddleware.MiddlewareFunc())
	{
		user.GET("/state", userAPI.GetState)
		user.POST("/logout", userAPI.Logout)
		user.GET("/devices", userAPI.GetDevices)
		user.DELETE("/devices/:id", userAPI.DeleteDevice)
	}

	profile := router.Group("profile")
	profile.Use(authMiddleware.MiddlewareFunc())
	{
		profile.GET("/", profileAPI.GetUserProfile)

		profile.POST("/country", profileAPI.UpdateUserCountry)
		profile.POST("/personal", profileAPI.UpdateUserPersonalInfo)
	}

	accounts := router.Group("accounts")
	accounts.Use(authMiddleware.MiddlewareFunc())
	{
		accounts.GET("/", accountsApi.GetAccounts)
		accounts.GET("/state/:id", accountsApi.GetAccountState)

		accounts.POST("/register", accountsApi.RegisterAccount)

		accounts.POST("/suspend/:id", accountsApi.SuspendAccount)
		accounts.POST("/resume/:id", accountsApi.ResumeAccount)
	}

	router.Run()
}
