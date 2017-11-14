package auth

import (
	"fmt"
	"github.com/gavrilaf/spawn/pkg/cache"
	"github.com/gavrilaf/spawn/pkg/env"
	mdl "github.com/gavrilaf/spawn/pkg/model"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	//"golang.org/x/net/context"
)

type StorageBridge struct {
	cache   *cache.RedisCache
	backend *pb.BackendBridge
}

func NewBridgeStorage(en *env.Environment) *StorageBridge {
	cache, err := cache.Connect(en)
	if err != nil {
		log.Errorf("Can not connect to cache: %v", err)
		return nil
	}
	log.Infof("Sessions Redis storage connected")

	backend, err := pb.CreateClient(en)
	if err != nil {
		log.Errorf("Can not connect to backend: %v", err)
		return nil
	}
	log.Infof("Connected to the backend")

	return &StorageBridge{cache, backend}
}

/////////////////////////////////////////////////////////////////////////////////

func (sb *StorageBridge) FindClient(id string) (mdl.Client, error) {
	return mdl.Client{}, fmt.Errorf("not implemented")
}

func (sb *StorageBridge) RegisterUser(username string, password string, device mdl.DeviceInfo) error {
	/*req := pb.RegisterUserRequest{
		Username:     username,
		PasswordHash: password,
		DeviceId:     deviceId}

	resp, err := sb.backend.Client.RegisterUser(context.Background(), &req)
	if resp != nil {
		log.Errorf("Registered user with id = %v", resp.Id)
	}*/

	return fmt.Errorf("not implemented")
}

func (sb *StorageBridge) FindUser(username string) (cache.AuthUser, error) {
	return cache.AuthUser{}, fmt.Errorf("not implemented")
}

func (sb *StorageBridge) FindDevice(userId string, deviceId string) (cache.AuthDevice, error) {
	return cache.AuthDevice{}, fmt.Errorf("not implemented")
}

func (sb *StorageBridge) AddDevice(userId string, device mdl.DeviceInfo) error {
	return fmt.Errorf("not implemented")
}

func (sb *StorageBridge) StoreSession(session cache.Session) error {
	return sb.cache.AddSession(session)
}

func (sb *StorageBridge) FindSession(id string) (cache.Session, error) {
	session, err := sb.cache.FindSession(id)
	if err != nil {
		log.Errorf("Can't find session with id %v, error: %v", id, err)
		return cache.Session{}, errSessionNotFound
	}

	return *session, nil
}
