package ginx

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func HandleAuthError(c *gin.Context, httpCode int, err error) {
	log.Errorf("API auth error: code=(%d), err=(%v)", httpCode, err)

	c.Header("WWW-Authenticate", "JWT realm="+defs.Realm)
	c.JSON(httpCode, gin.H{"error": errx.Error2Map(err, defs.ErrScope)})
	c.Abort()
}

func HandleError(c *gin.Context, scope string, httpCode int, err error) {
	log.Errorf("API error: code=(%d), scope=(%s), err=(%v)", httpCode, scope, err)

	errJSON := errx.Error2Map(err, scope)
	c.JSON(httpCode, gin.H{"error": errJSON})
	c.Abort()
}

func GetContextSession(c *gin.Context) (*mdl.Session, error) {
	session, exists := c.Get(defs.SessionKey)
	if !exists {
		return nil, defs.ErrSessionNotFound
	}
	return session.(*mdl.Session), nil
}
