package user

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	errScope = "user"
)

type Api interface {
	GetState(c *gin.Context)
	Logout(c *gin.Context)
}

///////////////////////////////////////////////////////////////////////////////

type ApiImpl struct {
	*api.Bridge
}

func CreateApi(bridge *api.Bridge) ApiImpl {
	return ApiImpl{Bridge: bridge}
}

///////////////////////////////////////////////////////////////////////////////

func (api ApiImpl) handleError(c *gin.Context, httpCode int, err error) {
	log.Errorf("profile.HandleError, code=%d, err=%v", httpCode, err)
	errJSON := errx.Error2Map(err, errScope)
	c.JSON(httpCode, gin.H{"error": errJSON})
	c.Abort()
}

func (api ApiImpl) getSession(c *gin.Context) (*mdl.Session, error) {
	return api.ReadModel.FindSession(c.GetString("session_id"))
}

///////////////////////////////////////////////////////////////////////////////

func (p ApiImpl) GetState(c *gin.Context) {
	session, err := p.getSession(c)
	if err != nil {
		log.Errorf("UserApi.GetState, could not find session: %v", err)
		p.handleError(c, http.StatusUnauthorized, err)
		return
	}

	state := UserState{
		UserID: session.UserID,
		Locale: session.Locale,
		Lang:   session.Lang,
		Permissions: auth.PermissionsDTO{
			IsDeviceConfirmed: session.IsDeviceConfirmed,
			Is2FARequired:     session.Is2FARequired,
			IsEmailConfirmed:  session.IsEmailConfirmed,
			IsLocked:          session.IsLocked,
			Scopes:            session.Scope,
		},
	}

	c.JSON(http.StatusOK, state.ToMap())
}

func (p ApiImpl) Logout(c *gin.Context) {
	session, err := p.getSession(c)
	if err != nil {
		log.Errorf("UserApi.GetState, could not find session: %v", err)
		p.handleError(c, http.StatusUnauthorized, err)
		return
	}

	err = p.ReadModel.DeleteSession(session.ID)
	if err != nil {
		log.Errorf("UserApi.GetState, could not invalidate session: %v", err)
		p.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, api.EmptySuccessResponse)
}
