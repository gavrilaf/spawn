package backend

import (
	//"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	mdl "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (srv *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.IDResponse, error) {
	log.Infof("CreateUser, %v", spew.Sdump(req))

	// Add user to the DB
	device := mdl.DeviceInfo{
		ID:          req.Device.ID,
		Name:        req.Device.Name,
		IsConfirmed: true, // Confirm device when user is registering
		Locale:      req.Device.Locale,
		Lang:        req.Device.Lang,
	}

	profile, err := srv.db.RegisterUser(req.Username, req.PasswordHash, device)
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

	log.Infof("User created, %v, %v", req.Username, profile.ID)

	return &pb.IDResponse{ID: profile.ID}, nil
}

func (srv *Server) AddDevice(ctx context.Context, req *pb.AddDeviceRequest) (*pb.Empty, error) {
	log.Infof("AddDevice, %v", spew.Sdump(req))

	device := mdl.DeviceInfo{
		ID:          req.Device.ID,
		UserID:      req.UserID,
		Name:        req.Device.Name,
		IsConfirmed: false,
		Locale:      req.Device.Locale,
		Lang:        req.Device.Lang,
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
	if err := srv.cache.AddConfirmCode("device", req.UserID+req.Device.ID, code); err != nil {
		log.Errorf("Storing confirm code error, %v", err)
		return nil, err
	}

	log.Infof("Device %v for user %v added. Confirm code %v", req.UserID, req.Device.ID, code)

	return &pb.Empty{}, nil
}

func (srv *Server) HandleLogin(ctx context.Context, req *pb.LoginRequest) (*pb.Empty, error) {
	log.Infof("HandleLogin, %v", spew.Sdump(req))

	// Add profile to the cache
	err := srv.updateCachedUserProfile(req.UserID)
	if err != nil {
		return nil, err
	}

	// Update device info (name & locale)
	device, err := srv.db.GetUserDevice(req.UserID, req.Device.ID)
	if err != nil {
		log.Errorf("Could not find device for (%v, %v): %v", req.UserID, req.Device.ID, err)
		return nil, err
	}

	device.Name = req.Device.Name
	device.Lang = req.Device.Lang
	device.Locale = req.Device.Locale

	if err = srv.db.UpdateDevice(*device); err != nil {
		log.Errorf("Could not update device %v: %v", spew.Sdump(device), err)
		// It isn't critical, continue execution
	}

	// Update last login info
	err = srv.db.LogUserLogin(req.UserID, req.Device.ID, req.UserAgent, req.LoginIP, req.LoginRegion)
	if err != nil {
		log.Errorf("Could not log login info, %v", err)
		return nil, err
	}

	return &pb.Empty{}, nil
}
