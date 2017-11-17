package auth

import (
	"context"

	"github.com/gavrilaf/spawn/pkg/cache"
	mdl "github.com/gavrilaf/spawn/pkg/model"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
)

func (sb StorageImpl) FindClient(id string) (mdl.Client, error) {
	var storageMock = NewStorageMock(nil)
	return storageMock.FindClient(id)
}

func (sb StorageImpl) RegisterUser(username string, password string, device mdl.DeviceInfo) error {

	req := pb.CreateUserRequest{
		Username:     username,
		PasswordHash: password,
		Device: &pb.Device{
			Id:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	resp, err := sb.Backend.Client.CreateUser(context.Background(), &req)
	if err == nil && resp != nil {
		log.Infof("Registered user (%v, %v)", username, resp.UserId)
	}

	return err
}

func (sb StorageImpl) FindUser(username string) (*cache.AuthUser, error) {
	return sb.Cache.FindUserAuthInfo(username)
}

func (sb StorageImpl) FindDevice(userId string, deviceId string) (*cache.AuthDevice, error) {
	return sb.Cache.FindDevice(userId, deviceId)
}

func (sb StorageImpl) AddDevice(userId string, device mdl.DeviceInfo) (*cache.AuthDevice, error) {

	req := pb.AddDeviceRequest{
		UserId: userId,
		Device: &pb.Device{
			Id:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	_, err := sb.Backend.Client.AddDevice(context.Background(), &req)
	if err == nil {
		log.Infof("Added device (%v, %v)", userId, device.ID)
	}

	return sb.FindDevice(userId, device.ID)
}

func (sb StorageImpl) StoreSession(session cache.Session) error {
	return sb.Cache.AddSession(session)
}

func (sb StorageImpl) FindSession(id string) (*cache.Session, error) {
	session, err := sb.Cache.FindSession(id)
	if err != nil {
		log.Errorf("Can't find session with id %v, error: %v", id, err)
		return nil, errSessionNotFound
	}

	return session, nil
}
