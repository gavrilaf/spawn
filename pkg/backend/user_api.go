package backend

import (
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	"github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func (self *Server) AddDevice(arg *pb.UserDevice) (*pb.Empty, error) {
	log.Infof("AddDevice, %s", arg.String())

	device := mdl.DeviceInfo{
		ID:          arg.Device.ID,
		UserID:      arg.UserID,
		Name:        arg.Device.Name,
		IsConfirmed: false,
		Locale:      arg.Device.Locale,
		Lang:        arg.Device.Lang,
	}

	if err := self.db.AddDevice(device); err != nil {
		log.Errorf("Could not add device to the db, %v", err)
		return nil, err
	}

	if err := self.cache.SetDevice(device); err != nil {
		log.Errorf("Could not add device to the cache, %v", err)
		return nil, err
	}

	// Generate confirm code
	code := cryptx.GenerateConfimCode()
	if err := self.cache.AddDeviceConfirmCode(arg.UserID, arg.Device.ID, code); err != nil {
		log.Errorf("Storing confirm code error, %v", err)
		return nil, err
	}

	log.Infof("Device %s for user %s added. Confirm code %s", arg.UserID, arg.Device.ID, code)

	return &pb.Empty{}, nil
}

func (self *Server) ConfirmDevice(arg *pb.ConfirmDeviceReq) (*pb.Empty, error) {
	log.Infof("ConfirmDevice: %s", arg.String())

	session, err := self.cache.GetSession(arg.SessionId)
	if err != nil {
		log.Errorf("Could not read session from the read model, %v", err)
		return nil, err
	}

	code, err := self.cache.GetDeviceConfirmCode(session.UserID, session.DeviceID)
	if err != nil {
		log.Errorf("Could not read confirm code from the read model, %v", err)
		return nil, err
	}

	if len(code) == 0 || code != arg.Code {
		log.Errorf("Could not find confirmation code")
		return nil, errx.ErrNotFound(ErrScope, "Could not find confirmation code")
	}

	err = self.db.ConfirmDevice(session.UserID, session.DeviceID)
	if err != nil {
		log.Errorf("Could not update device, %v", err)
		return nil, err
	}

	err = self.updateCachedUserDevices(session.UserID)
	if err != nil {
		log.Errorf("Could not update devices list in the read model, %v", err)
		return nil, err
	}

	session.IsDeviceConfirmed = true
	err = self.cache.SetSession(*session)
	if err != nil {
		log.Errorf("Could not update session in the read model, %v", err)
		return nil, err
	}

	err = self.cache.DeleteConfirmCode(session.UserID, session.DeviceID)
	if err != nil {
		log.Errorf("Could not delete confirm code from the read model, %v", err)
		return nil, err

	}

	log.Infof("Device (%s, %s) is confirmed with code %s", session.UserID, session.DeviceID, arg.Code)

	return &pb.Empty{}, nil
}

func (self *Server) DeleteDevice(arg *pb.UserDeviceID) (*pb.Empty, error) {
	log.Infof("DeleteDevice: %s", arg.String())

	err := self.db.RemoveDevice(arg.UserID, arg.DeviceID)
	if err != nil {
		log.Errorf("Could not delete device from the write model, %v", err)
		return nil, err
	}

	err = self.updateCachedUserDevices(arg.UserID)
	if err != nil {
		log.Errorf("Could not update devices list in the read model, %v", err)
		return nil, err
	}

	log.Infof("Device (%v, %v) is deleted", arg.UserID, arg.DeviceID)

	// Invalidate session
	session, _ := self.cache.FindSession(arg.UserID, arg.DeviceID)
	if session != nil {
		err = self.cache.DeleteSession(session.ID)
		if err != nil {
			log.Errorf("Could not invalidate session with id %s, %v", session.ID, err)
		} else {
			log.Infof("Session %s is invalidated", session.ID)
		}
	}

	return &pb.Empty{}, nil
}
