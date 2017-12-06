package auth

import (
	"context"

	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
)

func (sb StorageImpl) Close() {
	if sb.Bridge != nil {
		sb.Bridge.Close()
	}
}

func (sb StorageImpl) FindClient(id string) (*db.Client, error) {
	var storageMock = NewStorageMock(nil)
	return storageMock.FindClient(id)
}

func (sb StorageImpl) RegisterUser(username string, password string, device db.DeviceInfo) error {
	req := pb.CreateUserRequest{
		Username:     username,
		PasswordHash: password,
		Device: &pb.Device{
			ID:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	resp, err := sb.Backend.Client.CreateUser(context.Background(), &req)
	if err == nil && resp != nil {
		log.Infof("Registered user (%v, %v)", username, resp.ID)
	}

	return err
}

func (sb StorageImpl) FindUser(username string) (*mdl.AuthUser, error) {
	return sb.Cache.FindUserAuthInfo(username)
}

func (sb StorageImpl) FindDevice(userId string, deviceId string) (*mdl.AuthDevice, error) {
	return sb.Cache.FindDevice(userId, deviceId)
}

func (sb StorageImpl) AddDevice(userID string, device db.DeviceInfo) (*mdl.AuthDevice, error) {

	req := pb.AddDeviceRequest{
		UserID: userID,
		Device: &pb.Device{
			ID:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	_, err := sb.Backend.Client.AddDevice(context.Background(), &req)
	if err == nil {
		log.Infof("Added device (%v, %v)", userID, device.ID)
	}

	return sb.FindDevice(userID, device.ID)
}

func (sb StorageImpl) StoreSession(session mdl.Session) error {
	return sb.Cache.AddSession(session)
}

func (sb StorageImpl) FindSession(id string) (*mdl.Session, error) {
	session, err := sb.Cache.FindSession(id)
	if err != nil {
		log.Errorf("Can't find session with id %v, error: %v", id, err)
		return nil, errSessionNotFound
	}

	return session, nil
}

func (sb StorageImpl) HandlerLogin(session mdl.Session, ctx LoginContext) error {

	req := pb.LoginRequest{
		SessionID: session.ID,
		UserID:    session.UserID,
		Device: &pb.Device{
			ID:     session.DeviceID,
			Name:   ctx.DeviceName,
			Lang:   session.Lang,
			Locale: session.Locale},
		LoginIP:     ctx.IP,
		LoginRegion: ctx.Region}

	_, err := sb.Backend.Client.HandleLogin(context.Background(), &req)
	if err != nil {
		log.Errorf("Register login error, %v", err)
	}
	return err
}
