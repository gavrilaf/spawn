package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/api/utils"
)

type Access struct {
	*api.Bridge
	defAccess types.Access
	acl       map[string]types.Access
}

func CreateAccess(bridge *api.Bridge, defAccess types.Access, acl []types.EndpointAccess) Access {
	p := Access{Bridge: bridge, defAccess: defAccess, acl: make(map[string]types.Access)}
	for _, acc := range acl {
		p.acl[types.GetEndpointKey(acc.Group, acc.Endpoint)] = acc.Access
	}
	return p
}

func (self *Access) checkAccess(c *gin.Context) error {
	log.Infof("Check access for handler: %s", c.GetString(types.EndpointKey))
	return nil
}

func (self Access) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := self.checkAccess(c)
		if err != nil {
			utils.HandleAuthError(c, http.StatusForbidden, err)
		}
		return
	}
}
