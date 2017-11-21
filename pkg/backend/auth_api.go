package backend

import (
	//"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	mdl "github.com/gavrilaf/spawn/pkg/dbx/model"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (srv *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Infof("CreateUser, %v", spew.Sdump(req))

	// Add user to the DB
	device := mdl.DeviceInfo{
		ID:          req.Device.Id,
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

	return &pb.CreateUserResponse{UserId: profile.ID}, nil
}

func (srv *Server) AddDevice(ctx context.Context, req *pb.AddDeviceRequest) (*pb.Empty, error) {
	log.Infof("AddDevice, %v", spew.Sdump(req))

	device := mdl.DeviceInfo{
		ID:          req.Device.Id,
		Name:        req.Device.Name,
		IsConfirmed: false,
		Locale:      req.Device.Locale,
		Lang:        req.Device.Lang,
	}

	if err := srv.db.AddDevice(req.UserId, device); err != nil {
		log.Errorf("Could not add device to the db, %v", err)
		return nil, err
	}

	if err := srv.cache.SetDevice(req.UserId, device); err != nil {
		log.Errorf("Could not add device to the cache, %v", err)
		return nil, err
	}

	// Generate confirm code
	code := cryptx.GenerateConfimCode()
	if err := srv.cache.AddConfirmCode("device", req.UserId+req.Device.Id, code); err != nil {
		log.Errorf("Storing confirm code error, %v", err)
		return nil, err
	}

	log.Infof("Device %v for user %v added. Confirm code %v", req.UserId, req.Device.Id, code)

	return &pb.Empty{}, nil
}

func (srv *Server) HandleLogin(ctx context.Context, req *pb.LoginRequest) (*pb.Empty, error) {
	log.Infof("HandleLogin, %v", spew.Sdump(req))

	// Add profile to the cache

	return &pb.Empty{}, nil
}

func (srv *Server) DoConfirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.Empty, error) {
	log.Infof("DoConfirm, %v", spew.Sdump(req))

	switch req.Kind {

	}

	return &pb.Empty{}, nil
}
