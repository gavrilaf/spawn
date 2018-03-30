package user

import (
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"net/http"

	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/types"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (p ApiImpl) GetState(c *gin.Context) {
	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("Usertypes.GetState, could not find session: %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, types.ErrSessionNotFound)
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
		log.Errorf("Usertypes.GetState, could not find session: %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, types.ErrSessionNotFound)
		return
	}

	err = p.ReadModel.DeleteSession(session.ID)
	if err != nil {
		log.Errorf("Usertypes.GetState, could not invalidate session: %v", err)
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, types.EmptySuccessResponse)
}

///////////////////////////////////////////////////////////////////////////////

func (p ApiImpl) GetDevices(c *gin.Context) {
	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("Usertypes.GetDevices, could not find session: %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, types.ErrSessionNotFound)
		return
	}

	devices, err := p.ReadModel.GetUserDevicesInfo(session.UserID)
	if err != nil {
		log.Errorf("Usertypes.GetDevices, could not read devices: %v", err)
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, err)
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
		log.Errorf("Profiletypes.ConfirmDevice, could not bind, %v", err)
		p.HandleError(c, types.ErrScope, http.StatusBadRequest, err)
		return
	}

	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("Profiletypes.ConfirmDevice, could not find session, %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, err)
		return
	}

	if session.IsDeviceConfirmed {
		log.Errorf("Profiletypes.ConfirmDevice, device (%v, %v) already confirmed", session.UserID, session.DeviceID)
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, types.ErrAlreadyConfirmed)
		return
	}

	log.Infof("Profiletypes.ConfirmDevice, confirm device (%v, %v) with code %v", session.UserID, session.DeviceID, req.Code)

	_, err = p.WriteModel.DoConfirm(&pb.ConfirmDeviceReq{
		Code: req.Code,
		Kind: "device"})

	if err != nil {
		log.Errorf("Profiletypes.ConfirmDevice, confirm device (%v, %v) error %v", session.UserID, session.DeviceID, err)
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, err)
		return
	}

	log.Infof("Profiletypes.ConfirmDevice, device (%v, %v) is confirmed", session.UserID, session.DeviceID)

	c.JSON(200, types.EmptySuccessResponse)
}

func (p ApiImpl) DeleteDevice(c *gin.Context) {
	deviceID := c.Param("id")

	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("Usertypes.DeleteDevice, could not find session, %v", err)
		p.HandleError(c, types.ErrScope, http.StatusUnauthorized, err)
		return
	}

	if session.DeviceID == deviceID {
		log.Errorf("Usertypes.DeleteDevice, could not delete active device")
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, types.ErrDeleteCurrentDevice)
		return
	}

	_, err = p.Bridge.WriteModel.DeleteDevice(&pb.UserDeviceID{
		UserID:   session.UserID,
		DeviceID: deviceID})

	if err != nil {
		log.Errorf("Usertypes.DeleteDevice, could not delete device, %v", err)
		p.HandleError(c, types.ErrScope, http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, types.EmptySuccessResponse)
}
