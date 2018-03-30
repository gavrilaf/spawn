package access

import (
	"github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gin-gonic/gin"
)

type Checker struct {
	routes map[string]string // HandlerName -> Route map (gin specific)
}

func CreateChecker() *Checker {
	routes := make(map[string]string)
	return &Checker{routes: routes}
}

func (ac *Checker) Init(routes *gin.RoutesInfo, defAccess types.Access, acl []types.EndpointAccess) {

}

func (ac *Checker) CheckAccess(session *mdl.Session, c *gin.Context) error {
	return nil
}
