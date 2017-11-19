package profile

import (
	"net/http"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gin-gonic/gin"
	//"github.com/sirupsen/logrus"
)

func (api ProfileApiImpl) WhoAmI(c *gin.Context) {
	session, err := api.getSession(c)
	if err != nil {
		api.handleError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(200, gin.H{
		"client_id": session.ClientID,
		"user_id":   session.UserID,
		"device_id": session.DeviceID,
	})
}

func (pi ProfileApiImpl) ConfirmDevice(c *gin.Context) {
	var req ConfirmDeviceRequest

	err := c.Bind(&req)
	if err != nil {
		pi.handleError(c, http.StatusBadRequest, err)
		return
	}

	session, err := pi.getSession(c)
	if err != nil {
		pi.handleError(c, http.StatusUnauthorized, err)
		return
	}

	if session.IsDeviceConfirmed {
		pi.handleError(c, http.StatusUnauthorized, errAlreadyConfirmed)
		return
	}

	code, err := pi.Cache.GetConfirmCode("device", session.UserID+session.DeviceID)
	if err != nil {
		pi.handleError(c, http.StatusInternalServerError, err)
		return
	}

	if code != req.Code {
		pi.handleError(c, http.StatusInternalServerError, errInvalidConfirm)
		return
	}

	c.JSON(200, api.EmptySuccessResponse)
}
