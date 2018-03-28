package main

import (
	"os"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/config"
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

	log.Infof("REST service environment: %s", env.String())

	apiBridge := api.CreateBridge(env)
	if apiBridge == nil {
		log.Info("Could not connect to the api bridge")
	}

	config.ConfigureRouter(router, apiBridge)

	router.Run()
}
