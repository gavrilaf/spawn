package profile

import (
	"net/http"
	//"time"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetUserProfile - return current user profile
func (p ApiImpl) GetUserProfile(c *gin.Context) {
	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("ProfileApi.GetUserProfile, could not find session, %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	profile, err := p.ReadModel.GetUserProfile(session.UserID)
	if err != nil {
		log.Errorf("ProfileApi.GetUserProfile, read profile %v error, %v", session.UserID, err)
		p.HandleError(c, api.ErrScope, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, profile.ToMap())
}

// UpdateUserCountry - update current user Country
func (p ApiImpl) UpdateUserCountry(c *gin.Context) {
	var req UpdateCountryRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserCountry, could not bind, %v", err)
		p.HandleError(c, api.ErrScope, http.StatusBadRequest, api.ErrInvalidRequest)
		return
	}

	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserCountry, could not find session, %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	log.Infof("ProfileApi.UpdateUserCountry, for user %v country %v", session.UserID, req.Country)

	_, err = p.WriteModel.UpdateUserCountry(&pb.UserCountry{
		UserID:  session.UserID,
		Country: req.Country,
	})

	if err != nil {
		log.Errorf("ProfileApi.UpdateUserCountry, backend error: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, api.EmptySuccessResponse)
}

// UpdateUserPersonalInfo - update current user personal info
func (p ApiImpl) UpdateUserPersonalInfo(c *gin.Context) {
	var req UpdatePersonalInfoRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserPersonalInfo, could not bind: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserPersonalInfo, could not find session: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	log.Infof("ProfileApi.UpdateUserPersonalInfo, for user %s, %v", session.UserID, req)

	t := utils.ParseBirthdayDate(req.BirthDate)
	_, err = p.WriteModel.UpdateUserPersonalInfo(&pb.UserPersonalInfo{
		UserID:    session.UserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: &pb.BirthDate{Year: int32(t.Year()), Month: int32(t.Month()), Day: int32(t.Day())},
	})

	if err != nil {
		log.Errorf("ProfileApi.UpdateUserPersonalInfo, backend error: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, api.EmptySuccessResponse)
}
