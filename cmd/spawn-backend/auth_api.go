package main

import (
	//"fmt"
	"github.com/davecgh/go-spew/spew"
	mdl "github.com/gavrilaf/spawn/pkg/model"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (srv *BackendServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Infof("CreateUser, %v", spew.Sdump(req))

	// Add user to the DB
	device := mdl.DeviceInfo{
		ID:          req.Device.Id,
		Name:        req.Device.Name,
		IsConfirmed: true, // Confirm device when user is registering
		Locale:      req.Device.Locale,
		Lang:        req.Device.Lang,
	}

	profile, err := srv.Db.RegisterUser(req.Username, req.PasswordHash, device)
	if err != nil {
		log.Errorf("Could not add user to the db, %v", err)
		return nil, err
	}

	// Update Redis cache
	err = srv.Cache.AddUserAuthInfo(*profile, []mdl.DeviceInfo{device})
	if err != nil {
		log.Errorf("Could not add user to the cache, %v", err)
		return nil, err
	}

	// Send email

	log.Infof("User created, %v, %v", req.Username, profile.ID)

	return &pb.CreateUserResponse{UserId: profile.ID}, nil
}

func (srv *BackendServer) AddDevice(ctx context.Context, req *pb.AddDeviceRequest) (*pb.Empty, error) {
	log.Infof("AddDevice, %v", spew.Sdump(req))

	device := mdl.DeviceInfo{
		ID:     req.Device.Id,
		Name:   req.Device.Name,
		Locale: req.Device.Locale,
		Lang:   req.Device.Lang,
	}

	if err := srv.Db.AddDevice(req.UserId, device); err != nil {
		log.Errorf("Could not add device to the db, %v", err)
	}

	if err := srv.Cache.AddDevice(req.UserId, device); err != nil {
		log.Errorf("Could not add device to the cache, %v", err)
	}

	return &pb.Empty{}, nil
}

func (srv *BackendServer) RegisterLogin(ctx context.Context, req *pb.LoginRequest) (*pb.Empty, error) {
	log.Infof("RegisterLogin, %v", spew.Sdump(req))

	return &pb.Empty{}, nil
}
