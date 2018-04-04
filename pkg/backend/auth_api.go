package backend

import (
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	mdl "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	log "github.com/sirupsen/logrus"
)

func (srv *Server) CreateUser(arg *pb.CreateUserReq) (*pb.ResID, error) {
	log.Infof("CreateUser, %s", arg.String())

	// Add user to the DB
	device := mdl.DeviceInfo{
		ID:          arg.Device.ID,
		Name:        arg.Device.Name,
		IsConfirmed: true, // Confirm device when user is registering
		Locale:      arg.Device.Locale,
		Lang:        arg.Device.Lang,
	}

	profile, err := srv.db.RegisterUser(arg.Username, arg.PasswordHash, device)
	if err != nil {
		log.Errorf("Could not add user to the db, %v", err)
		return nil, err
	}

	// Update Redis cache
	err = srv.cache.SetUserAuthInfo(*profile, []mdl.DeviceInfo{device})
	if err != nil {
		log.Errorf("Could not add user to the cache, %v", err)
		return nil, err
	}

	// Send email

	log.Infof("User created, %s, %s", arg.Username, profile.ID)

	return &pb.ResID{ID: profile.ID}, nil
}

func (srv *Server) AddDevice(arg *pb.UserDevice) (*pb.Empty, error) {
	log.Infof("AddDevice, %s", arg.String())

	device := mdl.DeviceInfo{
		ID:          arg.Device.ID,
		UserID:      arg.UserID,
		Name:        arg.Device.Name,
		IsConfirmed: false,
		Locale:      arg.Device.Locale,
		Lang:        arg.Device.Lang,
	}

	if err := srv.db.AddDevice(device); err != nil {
		log.Errorf("Could not add device to the db, %v", err)
		return nil, err
	}

	if err := srv.cache.SetDevice(device); err != nil {
		log.Errorf("Could not add device to the cache, %v", err)
		return nil, err
	}

	// Generate confirm code
	code := cryptx.GenerateConfimCode()
	if err := srv.cache.AddDeviceConfirmCode(arg.UserID, arg.Device.ID, code); err != nil {
		log.Errorf("Storing confirm code error, %v", err)
		return nil, err
	}

	log.Infof("Device %s for user %s added. Confirm code %s", arg.UserID, arg.Device.ID, code)

	return &pb.Empty{}, nil
}

func (srv *Server) HandleLogin(arg *pb.LoginReq) (*pb.Empty, error) {
	log.Infof("HandleLogin, %s", arg.String())

	// Add profile to the cache
	err := srv.updateCachedUserProfile(arg.UserID)
	if err != nil {
		return nil, err
	}

	// Update device info (name & locale)
	device, err := srv.db.GetUserDevice(arg.UserID, arg.Device.ID)
	if err != nil {
		log.Errorf("Could not find device for (%s, %s): %v", arg.UserID, arg.Device.ID, err)
		return nil, err
	}

	device.Name = arg.Device.Name
	device.Lang = arg.Device.Lang
	device.Locale = arg.Device.Locale

	if err = srv.db.UpdateDevice(*device); err != nil {
		log.Errorf("Could not update device %v", err)
		// It isn't critical, continue execution
	}

	// Update last login info
	err = srv.db.LogUserLogin(arg.UserID, arg.Device.ID, arg.UserAgent, arg.LoginIP, arg.LoginRegion)
	if err != nil {
		log.Errorf("Could not log login info, %v", err)
		return nil, err
	}

	err = srv.updateCachedUserDevices(arg.UserID)
	if err != nil {
		log.Errorf("Could not update user devices, %v", err)
	}

	return &pb.Empty{}, nil
}
