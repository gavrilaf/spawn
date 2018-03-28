package auth

import (
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gin-gonic/gin"
)

type accessChecker struct {
	routes map[string]string // HandlerName -> Route map (gin specific)
}

func createAccessChecker() *accessChecker {
	routes := make(map[string]string)
	return &accessChecker{routes: routes}
}

func (ac *accessChecker) checkAccess(session *mdl.Session, c *gin.Context) error {
	return nil
}
