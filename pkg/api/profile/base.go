package profile

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gin-gonic/gin"
)

const (
	errScope = "profile"
)

type Api interface {
	GetUserProfile(c *gin.Context)
	UpdateUserCountry(c *gin.Context)
	UpdateUserPersonalInfo(c *gin.Context)
}

///////////////////////////////////////////////////////////////////////////////

type ApiImpl struct {
	*api.Bridge
}

func CreateApi(bridge *api.Bridge) Api {
	return ApiImpl{Bridge: bridge}
}
