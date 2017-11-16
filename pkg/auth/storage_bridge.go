package auth

import (
	//"fmt"
	"github.com/gavrilaf/spawn/pkg/cache"
	"github.com/gavrilaf/spawn/pkg/env"
	mdl "github.com/gavrilaf/spawn/pkg/model"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type StorageBridge struct {
	cache   *cache.Cache
	backend *pb.BackendBridge
}

func NewBridgeStorage(en *env.Environment) *StorageBridge {
	log.Infof("Starting storage with environment: %v", en.GetName())

	cache, err := cache.Connect(en)
	if err != nil {
		log.Errorf("Can not connect to cache: %v", err)
		return nil
	}
	log.Infof("Cache connection, ok")

	backend, err := pb.CreateClient(en)
	if err != nil {
		log.Errorf("Can not connect to backend: %v", err)
		return nil
	}
	log.Infof("Backend connection, ok")

	return &StorageBridge{cache, backend}
}

/////////////////////////////////////////////////////////////////////////////////

var storageMock = NewStorageMock(nil)

func (sb *StorageBridge) FindClient(id string) (mdl.Client, error) {
	return storageMock.FindClient(id)
}

func (sb *StorageBridge) RegisterUser(username string, password string, device mdl.DeviceInfo) error {

	req := pb.CreateUserRequest{
		Username:     username,
		PasswordHash: password,
		Device: &pb.Device{
			Id:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	resp, err := sb.backend.Client.CreateUser(context.Background(), &req)
	if err == nil && resp != nil {
		log.Infof("Registered user (%v, %v)", username, resp.UserId)
	}

	return err
}

func (sb *StorageBridge) FindUser(username string) (*cache.AuthUser, error) {
	return sb.cache.FindUserAuthInfo(username)
}

func (sb *StorageBridge) FindDevice(userId string, deviceId string) (*cache.AuthDevice, error) {
	return sb.cache.FindDevice(userId, deviceId)
}

func (sb *StorageBridge) AddDevice(userId string, device mdl.DeviceInfo) (*cache.AuthDevice, error) {

	req := pb.AddDeviceRequest{
		UserId: userId,
		Device: &pb.Device{
			Id:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	_, err := sb.backend.Client.AddDevice(context.Background(), &req)
	if err == nil {
		log.Infof("Added device (%v, %v)", userId, device.ID)
	}

	return sb.FindDevice(userId, device.ID)
}

func (sb *StorageBridge) StoreSession(session cache.Session) error {
	return sb.cache.AddSession(session)
}

func (sb *StorageBridge) FindSession(id string) (*cache.Session, error) {
	session, err := sb.cache.FindSession(id)
	if err != nil {
		log.Errorf("Can't find session with id %v, error: %v", id, err)
		return nil, errSessionNotFound
	}

	return session, nil
}
