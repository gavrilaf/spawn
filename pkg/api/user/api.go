package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/api/ginx"
	"github.com/gavrilaf/spawn/pkg/backend/pb"
)

func (self ApiImpl) GetState(c *gin.Context) {
	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("User.GetState, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
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

func (self ApiImpl) Logout(c *gin.Context) {
	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("User.GetState, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, defs.ErrSessionNotFound)
		return
	}

	err = self.ReadModel.DeleteSession(session.ID)
	if err != nil {
		log.Errorf("User.GetState, could not invalidate session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, defs.EmptySuccessResponse)
}

///////////////////////////////////////////////////////////////////////////////

func (self ApiImpl) GetDevices(c *gin.Context) {
	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("User.GetDevices, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	devices, err := self.ReadModel.GetUserDevicesInfo(session.UserID)
	if err != nil {
		log.Errorf("User.GetDevices, could not read devices: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
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

func (self ApiImpl) ConfirmDevice(c *gin.Context) {
	var req ConfirmDeviceCode

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("User.ConfirmDevice, could not bind, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusBadRequest, err)
		return
	}

	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("User.ConfirmDevice, could not find session, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	if session.IsDeviceConfirmed {
		log.Errorf("User.User, device (%s, %s) already confirmed", session.UserID, session.DeviceID)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, defs.ErrAlreadyConfirmed)
		return
	}

	log.Infof("User.ConfirmDevice, confirm device (%s, %s) with code %s", session.UserID, session.DeviceID, req.Code)

	_, err = self.WriteModel.ConfirmDevice(&pb.ConfirmDeviceReq{
		SessionId: session.ID,
		Code:      req.Code})

	if err != nil {
		log.Errorf("User.ConfirmDevice, confirm device (%s, %s) error %v", session.UserID, session.DeviceID, err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
		return
	}

	log.Infof("User.ConfirmDevice, device (%s, %s) is confirmed", session.UserID, session.DeviceID)
	c.JSON(200, defs.EmptySuccessResponse)
}

func (self ApiImpl) GetConfirmCode(c *gin.Context) {
	var req GetConfirmCodeRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("User.GetConfirmCode, could not bind, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusBadRequest, defs.ErrInvalidRequest)
		return
	}

	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("User.GetConfirmCode, could not find session, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	code, error := self.ReadModel.GetDeviceConfirmCode(session.UserID, req.DeviceID)
	if error != nil {
		log.Errorf("User.GetConfirmCode, could not find code, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
	}

	c.JSON(200, ConfirmDeviceCode{Code: code}.ToMap())
}

func (self ApiImpl) DeleteDevice(c *gin.Context) {
	deviceID := c.Param("id")

	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("User.DeleteDevice, could not find session, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	if session.DeviceID == deviceID {
		log.Errorf("User.DeleteDevice, could not delete active device")
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, defs.ErrDeleteCurrentDevice)
		return
	}

	_, err = self.Bridge.WriteModel.DeleteDevice(&pb.UserDeviceID{
		UserID:   session.UserID,
		DeviceID: deviceID})

	if err != nil {
		log.Errorf("User.DeleteDevice, could not delete device, %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
		return
	}

	log.Infof("User.DeleteDevice, device (%s, %s) is delete", session.UserID, deviceID)

	c.JSON(200, defs.EmptySuccessResponse)
}
