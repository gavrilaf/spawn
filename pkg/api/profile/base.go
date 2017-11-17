package profile

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gin-gonic/gin"
)

type ProfileApiImpl struct {
	*api.StorageBridge
}

type ProfileApi interface {
	WhoAmI(c *gin.Context)
	ConfirmDevice(c *gin.Context)
}
