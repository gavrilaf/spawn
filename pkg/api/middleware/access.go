package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/api/ginx"
)

type Access struct {
	*api.Bridge
	defAccess defs.Access
	acl       map[string]defs.Access
}

func CreateAccess(bridge *api.Bridge, defAccess defs.Access, acl []defs.EndpointAccess) Access {
	p := Access{Bridge: bridge, defAccess: defAccess, acl: make(map[string]defs.Access)}
	for _, acc := range acl {
		p.acl[defs.GetEndpointKey(acc.Group, acc.Endpoint)] = acc.Access
	}
	return p
}

func (self *Access) checkAccess(c *gin.Context) error {
	log.Infof("Check access for handler: %s", c.GetString(defs.EndpointKey))
	return nil
}

func (self Access) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := self.checkAccess(c)
		if err != nil {
			ginx.HandleAuthError(c, http.StatusForbidden, err)
		}
		return
	}
}
