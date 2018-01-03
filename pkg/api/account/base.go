package account

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gin-gonic/gin"
)

type Api interface {
	GetAccounts(c *gin.Context)
	GetAccountState(c *gin.Context)

	RegisterAccount(c *gin.Context)

	SuspendAccount(c *gin.Context)
	ResumeAccount(c *gin.Context)
}

///////////////////////////////////////////////////////////////////////////////

type ApiImpl struct {
	*api.Bridge
}

func CreateApi(bridge *api.Bridge) ApiImpl {
	return ApiImpl{Bridge: bridge}
}
