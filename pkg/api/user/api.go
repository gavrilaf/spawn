package user

import (
	"context"
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (p ApiImpl) GetState(c *gin.Context) {
	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("UserApi.GetState, could not find session: %v", err)
		p.HandleError(c, errScope, http.StatusUnauthorized, err)
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
	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("UserApi.GetState, could not find session: %v", err)
		p.HandleError(c, errScope, http.StatusUnauthorized, err)
		return
	}

	err = p.ReadModel.DeleteSession(session.ID)
	if err != nil {
		log.Errorf("UserApi.GetState, could not invalidate session: %v", err)
		p.HandleError(c, errScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, api.EmptySuccessResponse)
}

///////////////////////////////////////////////////////////////////////////////

func (p ApiImpl) GetDevices(c *gin.Context) {
	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("UserApi.GetDevices, could not find session: %v", err)
		p.HandleError(c, errScope, http.StatusUnauthorized, err)
		return
	}

	devices, err := p.ReadModel.GetUserDevicesInfo(session.UserID)
	if err != nil {
		log.Errorf("UserApi.GetDevices, could not read devices: %v", err)
		p.HandleError(c, errScope, http.StatusInternalServerError, err)
		return
	}

	for indx, _ := range devices {
		if devices[indx].ID == session.DeviceID {
			devices[indx].IsCurrent = true
		}
	}

	d := UserDevices{Devices: devices}

	c.JSON(http.StatusOK, d.ToMap())
}

func (p ApiImpl) ConfirmDevice(c *gin.Context) {
	var req ConfirmDeviceRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("ProfileApi.ConfirmDevice, could not bind, %v", err)
		p.HandleError(c, errScope, http.StatusBadRequest, err)
		return
	}

	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("ProfileApi.ConfirmDevice, could not find session, %v", err)
		p.HandleError(c, errScope, http.StatusUnauthorized, err)
		return
	}

	if session.IsDeviceConfirmed {
		log.Errorf("ProfileApi.ConfirmDevice, device (%v, %v) already confirmed", session.UserID, session.DeviceID)
		p.HandleError(c, errScope, http.StatusInternalServerError, errAlreadyConfirmed)
		return
	}

	log.Infof("ProfileApi.ConfirmDevice, confirm device (%v, %v) with code %v", session.UserID, session.DeviceID, req.Code)

	_, err = p.WriteModel.Client.DoConfirm(context.Background(), &pb.ConfirmRequest{
		Code: req.Code,
		Kind: "device"})

	if err != nil {
		log.Errorf("ProfileApi.ConfirmDevice, confirm device (%v, %v) error %v", session.UserID, session.DeviceID, err)
		p.HandleError(c, errScope, http.StatusInternalServerError, err)
		return
	}

	log.Infof("ProfileApi.ConfirmDevice, device (%v, %v) is confirmed", session.UserID, session.DeviceID)

	c.JSON(200, api.EmptySuccessResponse)
}
