package auth

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gin-gonic/gin"
)

type Api interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)

	RefreshToken(c *gin.Context)
}

///////////////////////////////////////////////////////////////////////////////

type ApiImpl struct {
	*api.Bridge
}

func CreateApi(bridge *api.Bridge) ApiImpl {
	return ApiImpl{Bridge: bridge}
}
