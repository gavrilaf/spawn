package profile

import (
	"context"
	"net/http"

	"github.com/gavrilaf/spawn/pkg/api"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (p ProfileApiImpl) ConfirmDevice(c *gin.Context) {
	var req ConfirmDeviceRequest

	err := c.Bind(&req)
	if err != nil {
		log.Errorf("ProfileApi.ConfirmDevice, could not bind, %v", err)
		p.handleError(c, http.StatusBadRequest, err)
		return
	}

	session, err := p.getSession(c)
	if err != nil {
		log.Errorf("ProfileApi.ConfirmDevice, could not find session, %v", err)
		p.handleError(c, http.StatusUnauthorized, err)
		return
	}

	if session.IsDeviceConfirmed {
		log.Errorf("ProfileApi.ConfirmDevice, device (%v, %v) already confirmed", session.UserID, session.DeviceID)
		p.handleError(c, http.StatusInternalServerError, errAlreadyConfirmed)
		return
	}

	log.Infof("ProfileApi.ConfirmDevice, confirm device (%v, %v) with code %v", session.UserID, session.DeviceID, req.Code)

	_, err = p.WriteModel.Client.DoConfirm(context.Background(), &pb.ConfirmRequest{
		Code: req.Code,
		Kind: "device"})

	if err != nil {
		log.Errorf("ProfileApi.ConfirmDevice, confirm device (%v, %v) error %v", session.UserID, session.DeviceID, err)
		p.handleError(c, http.StatusInternalServerError, err)
		return
	}

	log.Infof("ProfileApi.ConfirmDevice, device (%v, %v) is confirmed", session.UserID, session.DeviceID)

	c.JSON(200, api.EmptySuccessResponse)
}
