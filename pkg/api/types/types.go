package types

import (
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	Path   string
	Method string
}

type Access struct {
	NeedDevice bool
	NeedEmail  bool
	MinScope   int
}

type EndpointAccess struct {
	Group string
	Endpoint
	Access
}

type AccessChecker interface {
	Init(routes *gin.RoutesInfo, defAccess Access, acl []EndpointAccess)
	CheckAccess(session *mdl.Session, c *gin.Context) error
}

var EmptySuccessResponse = map[string]interface{}{"success": true}
