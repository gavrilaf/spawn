package profile

import (
	"github.com/gavrilaf/spawn/pkg/api"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ProfileApi interface {
	ConfirmDevice(c *gin.Context)
	GetUserProfile(c *gin.Context)
	UpdateUserCountry(c *gin.Context)
	UpdateUserPersonalInfo(c *gin.Context)
}

///////////////////////////////////////////////////////////////////////////////

type ProfileApiImpl struct {
	*api.Bridge
}

func CreateApi(bridge *api.Bridge) ProfileApi {
	return ProfileApiImpl{Bridge: bridge}
}

///////////////////////////////////////////////////////////////////////////////

func (api ProfileApiImpl) handleError(c *gin.Context, httpCode int, err error) {
	log.Errorf("profile.HandleError, code=%d, err=%v", httpCode, err)
	errJSON := errx.Error2Json(err, errScope)
	c.JSON(httpCode, gin.H{"error": errJSON})
	c.Abort()
}

func (api ProfileApiImpl) getSession(c *gin.Context) (*mdl.Session, error) {
	return api.Cache.FindSession(c.GetString("session_id"))
}

/*func (api ProfileApiImpl) getProfile(c *gin.Context) (*mdl.UserProfile, error) {
	userId := c.GetString("user_id")

	api.Cache
}*/
