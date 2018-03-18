package backend

import (
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	log "github.com/sirupsen/logrus"
)

func (srv *Server) DoConfirm(arg *pb.ConfirmDeviceReq) (*pb.Empty, error) {
	log.Infof("DoConfirm: %s", arg.String())
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
