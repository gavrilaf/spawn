package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (p ProfileApiImpl) GetUserProfile(c *gin.Context) {
	session, err := p.getSession(c)
	if err != nil {
		log.Errorf("ProfileApi.GetUserProfile, could not find session, %v", err)
		p.handleError(c, http.StatusUnauthorized, err)
		return
	}

	profile, err := p.Cache.GetUserProfile(session.UserID)
	if err != nil {
		log.Errorf("ProfileApi.GetUserProfile, read profile %v error, %v", session.UserID, err)
		p.handleError(c, http.StatusInternalServerError, err)
	}

	//log.Infof("Profile: %v", spew.Sdump(profile))
	c.JSON(http.StatusOK, profile.ToMap())
}

func (pi ProfileApiImpl) UpdateUserCountry(c *gin.Context) {
}

func (pi ProfileApiImpl) UpdateUserPersonalInfo(c *gin.Context) {
}
