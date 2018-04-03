package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/api/ginx"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
)

type Access struct {
	*api.Bridge
	defAccess defs.Access
	acl       map[string]defs.Access
}

func CreateAccess(bridge *api.Bridge, defAccess defs.Access, acl []defs.EndpointAccess) Access {
	p := Access{Bridge: bridge, defAccess: defAccess, acl: make(map[string]defs.Access)}
	for _, acc := range acl {
		endpointKey := defs.GetEndpointKey("/"+acc.Group, acc.Endpoint)
		log.Infof("Acl for: %s, %v", endpointKey, acc.Access)
		p.acl[endpointKey] = acc.Access
	}
	return p
}

func (self Access) checkAccess(c *gin.Context) error {
	endpoint := c.GetString(defs.EndpointKey)
	log.Infof("Check access for endpoint: %s", endpoint)

	session, err := ginx.GetContextSession(c)
	if err != nil {
		return err
	}

	if acl, found := self.acl[endpoint]; found {
		log.Infof("Found acl for: %s, %v", endpoint, acl)
		return isAccessAllowed(session, acl)
	} else {
		log.Infof("Use default acl for: %s", endpoint)
		return isAccessAllowed(session, self.defAccess)
	}
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

func isAccessAllowed(session *mdl.Session, endpointCfg defs.Access) error {
	if session.IsLocked && !endpointCfg.Locked {
		return defs.ErrUserLocked
	}

	if !session.IsDeviceConfirmed && endpointCfg.Device {
		return defs.ErrDeviceNotConfirmed
	}

	if session.IsEmailConfirmed && endpointCfg.Email {
		return defs.ErrEmailNotConfirmed
	}

	if int(session.Scope) < endpointCfg.Scope {
		return defs.ErrAccessForbiden
	}

	return nil
}
