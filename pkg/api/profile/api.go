package profile

import (
	"net/http"
	//"time"

	types "github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetUserProfile - return current user profile
func (p ApiImpl) GetUserProfile(c *gin.Context) {
	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("Profiletypes.GetUserProfile, could not find session, %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, types.ErrSessionNotFound)
		return
	}

	profile, err := p.ReadModel.GetUserProfile(session.UserID)
	if err != nil {
		log.Errorf("Profiletypes.GetUserProfile, read profile %v error, %v", session.UserID, err)
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, profile.ToMap())
}

// UpdateUserCountry - update current user Country
func (p ApiImpl) UpdateUserCountry(c *gin.Context) {
	var req UpdateCountryRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("Profiletypes.UpdateUserCountry, could not bind, %v", err)
		p.HandleError(c, types.ErrScope, http.StatusBadRequest, types.ErrInvalidRequest)
		return
	}

	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("Profiletypes.UpdateUserCountry, could not find session, %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, types.ErrSessionNotFound)
		return
	}

	log.Infof("Profiletypes.UpdateUserCountry, for user %v country %v", session.UserID, req.Country)

	_, err = p.WriteModel.UpdateUserCountry(&pb.UserCountry{
		UserID:  session.UserID,
		Country: req.Country,
	})

	if err != nil {
		log.Errorf("Profiletypes.UpdateUserCountry, backend error: %v", err)
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, types.EmptySuccessResponse)
}

// UpdateUserPersonalInfo - update current user personal info
func (p ApiImpl) UpdateUserPersonalInfo(c *gin.Context) {
	var req UpdatePersonalInfoRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("Profiletypes.UpdateUserPersonalInfo, could not bind: %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, types.ErrSessionNotFound)
		return
	}

	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("Profiletypes.UpdateUserPersonalInfo, could not find session: %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, types.ErrSessionNotFound)
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
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, types.EmptySuccessResponse)
}
