package utils

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func HandleAuthError(c *gin.Context, httpCode int, err error) {
	c.Header("WWW-Authenticate", "JWT realm="+types.Realm)
	log.Errorf("auth error, code=%d, err=%v", httpCode, err)
	c.JSON(httpCode, gin.H{"error": errx.Error2Map(err, types.ErrScope)})
	c.Abort()
}
