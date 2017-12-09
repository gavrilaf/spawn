package profile

import (
	"context"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gavrilaf/spawn/pkg/api"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
)

// GetUserProfile - return current user profile
func (p ProfileApiImpl) GetUserProfile(c *gin.Context) {
	session, err := p.getSession(c)
	if err != nil {
		log.Errorf("ProfileApi.GetUserProfile, could not find session, %v", err)
		p.handleError(c, http.StatusUnauthorized, err)
		return
	}

	profile, err := p.ReadModel.GetUserProfile(session.UserID)
	if err != nil {
		log.Errorf("ProfileApi.GetUserProfile, read profile %v error, %v", session.UserID, err)
		p.handleError(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, profile.ToMap())
}

// UpdateUserCountry - update current user Country
func (p ProfileApiImpl) UpdateUserCountry(c *gin.Context) {
	var req UpdateCountryRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserCountry, could not bind, %v", err)
		p.handleError(c, http.StatusBadRequest, err)
		return
	}

	session, err := p.getSession(c)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserCountry, could not find session, %v", err)
		p.handleError(c, http.StatusUnauthorized, err)
		return
	}

	log.Infof("ProfileApi.UpdateUserCountry, for user %v country %v", session.UserID, req.Country)

	_, err = p.WriteModel.Client.UpdateUserCountry(context.Background(), &pb.UserCountryRequest{
		UserID:  session.UserID,
		Country: req.Country,
	})

	if err != nil {
		log.Errorf("ProfileApi.UpdateUserCountry, backend error: %v", err)
		p.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, api.EmptySuccessResponse)
}

// UpdateUserPersonalInfo - update current user personal info
func (p ProfileApiImpl) UpdateUserPersonalInfo(c *gin.Context) {
	var req UpdatePersonalInfoRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserPersonalInfo, could not bind: %v", err)
		p.handleError(c, http.StatusBadRequest, err)
		return
	}

	session, err := p.getSession(c)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserPersonalInfo, could not find session: %v", err)
		p.handleError(c, http.StatusUnauthorized, err)
		return
	}

	log.Infof("ProfileApi.UpdateUserPersonalInfo, for user %v: %v", session.UserID, spew.Sdump(req))

	t, err := time.Parse(time.RFC3339, req.BirthDate)
	if err != nil {
		log.Errorf("ProfileApi.UpdateUserPersonalInfo, invalid BirthDate: %v", err)
		p.handleError(c, http.StatusBadRequest, err)
		return
	}

	protoTime, _ := ptypes.TimestampProto(t)

	_, err = p.WriteModel.Client.UpdateUserPersonalInfo(context.Background(), &pb.UserPersonalInfoRequest{
		UserID:    session.UserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: protoTime,
	})

	if err != nil {
		log.Errorf("ProfileApi.UpdateUserPersonalInfo, backend error: %v", err)
		p.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, api.EmptySuccessResponse)
}
