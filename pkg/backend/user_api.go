package backend

import (
	"github.com/davecgh/go-spew/spew"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (srv *Server) DoConfirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.Empty, error) {
	log.Infof("DoConfirm, %v", spew.Sdump(req))

	switch req.Kind {

	}

	return &pb.Empty{}, nil
}

func (srv *Server) DeleteDevice(ctx context.Context, req *pb.DeleteDeviceRequest) (*pb.Empty, error) {
	log.Infof("DeleteDevice, %v", spew.Sdump(req))

	err := srv.db.RemoveDevice(req.UserID, req.DeviceID)
	if err != nil {
		log.Errorf("Could not delete device from the write model, %v", err)
		return nil, err
	}

	err = srv.updateCachedUserDevices(req.UserID)
	if err != nil {
		log.Errorf("Could not update devices list in the read model, %v", err)
		return nil, err
	}

	log.Infof("Device (%v, %v) is deleted", req.UserID, req.DeviceID)

	// Invalidate session
	session, _ := srv.cache.FindSession(req.UserID, req.DeviceID)
	if session != nil {
		err = srv.cache.DeleteSession(session.ID)
		if err != nil {
			log.Errorf("Could not invalidate session with id %v, %v", session.ID, err)
		} else {
			log.Infof("Session %v is invalidated", session.ID)
		}
	}

	return &pb.Empty{}, nil
}