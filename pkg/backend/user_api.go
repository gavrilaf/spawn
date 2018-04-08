package backend

import (
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func (srv *Server) ConfirmDevice(arg *pb.ConfirmDeviceReq) (*pb.Empty, error) {
	log.Infof("ConfirmDevice: %s", arg.String())

	session, err := srv.cache.GetSession(arg.SessionId)
	if err != nil {
		log.Errorf("Could not read session from the read model, %v", err)
		return nil, err
	}

	code, err := srv.cache.GetDeviceConfirmCode(session.UserID, session.DeviceID)
	if err != nil {
		log.Errorf("Could not read confirm code from the read model, %v", err)
		return nil, err
	}

	if len(code) == 0 || code != arg.Code {
		return nil, errx.ErrNotFound(ErrScope, "Could not find confirmation code")
	}

	device, err := srv.db.GetUserDevice(session.UserID, session.DeviceID)
	if err != nil {
		log.Errorf("Could not find device, %v", err)
		return nil, err
	}

	device.IsConfirmed = true
	err = srv.db.UpdateDevice(*device)
	if err != nil {
		log.Errorf("Could not update device, %v", err)
		return nil, err
	}

	err = srv.updateCachedUserDevices(session.UserID)
	if err != nil {
		log.Errorf("Could not update devices list in the read model, %v", err)
		return nil, err
	}

	session.IsDeviceConfirmed = true
	err = srv.cache.SetSession(*session)
	if err != nil {
		log.Errorf("Could not update session in the read model, %v", err)
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (srv *Server) DeleteDevice(arg *pb.UserDeviceID) (*pb.Empty, error) {
	log.Infof("DeleteDevice: %s", arg.String())

	err := srv.db.RemoveDevice(arg.UserID, arg.DeviceID)
	if err != nil {
		log.Errorf("Could not delete device from the write model, %v", err)
		return nil, err
	}

	err = srv.updateCachedUserDevices(arg.UserID)
	if err != nil {
		log.Errorf("Could not update devices list in the read model, %v", err)
		return nil, err
	}

	log.Infof("Device (%v, %v) is deleted", arg.UserID, arg.DeviceID)

	// Invalidate session
	session, _ := srv.cache.FindSession(arg.UserID, arg.DeviceID)
	if session != nil {
		err = srv.cache.DeleteSession(session.ID)
		if err != nil {
			log.Errorf("Could not invalidate session with id %s, %v", session.ID, err)
		} else {
			log.Infof("Session %s is invalidated", session.ID)
		}
	}

	return &pb.Empty{}, nil
}
