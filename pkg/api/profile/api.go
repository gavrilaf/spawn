package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/api/ginx"
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/utils"
)

// GetUserProfile - return current user profile
func (p ApiImpl) GetUserProfile(c *gin.Context) {
	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("Profiletypes.GetUserProfile, could not find session, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	profile, err := p.ReadModel.GetUserProfile(session.UserID)
	if err != nil {
		log.Errorf("Profiletypes.GetUserProfile, read profile %v error, %v", session.UserID, err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, profile.ToMap())
}

// UpdateUserCountry - update current user Country
func (p ApiImpl) UpdateUserCountry(c *gin.Context) {
	var req UpdateCountryRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("Profiletypes.UpdateUserCountry, could not bind, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusBadRequest, defs.ErrInvalidRequest)
		return
	}

	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("Profiletypes.UpdateUserCountry, could not find session, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	log.Infof("Profiletypes.UpdateUserCountry, for user %s country %s", session.UserID, req.Country)

	_, err = p.WriteModel.UpdateUserCountry(&pb.UserCountry{
		UserID:  session.UserID,
		Country: req.Country,
	})

	if err != nil {
		log.Errorf("Profiletypes.UpdateUserCountry, backend error: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, defs.EmptySuccessResponse)
}

// UpdateUserPersonalInfo - update current user personal info
func (p ApiImpl) UpdateUserPersonalInfo(c *gin.Context) {
	var req UpdatePersonalInfoRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("Profiletypes.UpdateUserPersonalInfo, could not bind: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, defs.ErrSessionNotFound)
		return
	}

	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("Profiletypes.UpdateUserPersonalInfo, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, defs.ErrSessionNotFound)
		return
	}

	log.Infof("Profiletypes.UpdateUserPersonalInfo, for user %s, %v", session.UserID, req)

	t := utils.ParseBirthdayDate(req.BirthDate)
	_, err = p.WriteModel.UpdateUserPersonalInfo(&pb.UserPersonalInfo{
		UserID:    session.UserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: &pb.BirthDate{Year: int32(t.Year()), Month: int32(t.Month()), Day: int32(t.Day())},
	})

	if err != nil {
		log.Errorf("Profiletypes.UpdateUserPersonalInfo, backend error: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, defs.EmptySuccessResponse)
}
