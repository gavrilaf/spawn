package user

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gin-gonic/gin"
)

const (
	errScope = "user"
)

var errAlreadyConfirmed = errx.New(errScope, "device-already-confirmed")

type Api interface {
	GetState(c *gin.Context)
	Logout(c *gin.Context)

	ConfirmDevice(c *gin.Context)
	GetDevices(c *gin.Context)
}

///////////////////////////////////////////////////////////////////////////////

type ApiImpl struct {
	*api.Bridge
}

func CreateApi(bridge *api.Bridge) ApiImpl {
	return ApiImpl{Bridge: bridge}
}
