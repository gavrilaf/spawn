package user

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gin-gonic/gin"
)

type Api interface {
	GetState(c *gin.Context)
	Logout(c *gin.Context)

	ConfirmDevice(c *gin.Context)
	GetDevices(c *gin.Context)
	DeleteDevice(c *gin.Context)
}

///////////////////////////////////////////////////////////////////////////////

type ApiImpl struct {
	*api.Bridge
}

func CreateApi(bridge *api.Bridge) ApiImpl {
	return ApiImpl{Bridge: bridge}
}
