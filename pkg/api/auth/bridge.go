package auth

import (
	"context"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
)

type Bridge struct {
	*api.Bridge
}

func (sb Bridge) Close() {
	if sb.Bridge != nil {
		sb.Bridge.Close()
	}
}

func (sb Bridge) GetClient(id string) (*db.Client, error) {
	return sb.ReadModel.FindClient(id)
}

func (sb Bridge) RegisterUser(username string, password string, device db.DeviceInfo) error {
	req := pb.CreateUserRequest{
		Username:     username,
		PasswordHash: password,
		Device: &pb.Device{
			ID:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	resp, err := sb.WriteModel.Client.CreateUser(context.Background(), &req)
	if err == nil && resp != nil {
		log.Infof("Registered user (%v, %v)", username, resp.ID)
	}

	return err
}

func (sb Bridge) FindUser(username string) (*mdl.AuthUser, error) {
	return sb.ReadModel.FindUserAuthInfo(username)
}

func (sb Bridge) GetDevice(userId string, deviceId string) (*mdl.AuthDevice, error) {
	return sb.ReadModel.GetDevice(userId, deviceId)
}

func (sb Bridge) AddDevice(userID string, device db.DeviceInfo) (*mdl.AuthDevice, error) {

	req := pb.AddDeviceRequest{
		UserID: userID,
		Device: &pb.Device{
			ID:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	_, err := sb.WriteModel.Client.AddDevice(context.Background(), &req)
	if err == nil {
		log.Infof("Added device (%v, %v)", userID, device.ID)
	}

	return sb.GetDevice(userID, device.ID)
}

func (sb Bridge) AddSession(session mdl.Session) (string, error) {
	return sb.ReadModel.AddSession(session, false)
}

func (sb Bridge) UpdateSession(session mdl.Session) error {
	return sb.ReadModel.SetSession(session)
}

func (sb Bridge) GetSession(id string) (*mdl.Session, error) {
	session, err := sb.ReadModel.GetSession(id)
	if err != nil {
		log.Errorf("Can't find session with id %v, error: %v", id, err)
		return nil, api.ErrSessionNotFound
	}

	return session, nil
}

func (sb Bridge) HandlerLogin(session mdl.Session, ctx LoginContext) error {

	req := pb.LoginRequest{
		SessionID: session.ID,
		UserID:    session.UserID,
		Device: &pb.Device{
			ID:     session.DeviceID,
			Name:   ctx.DeviceName,
			Lang:   session.Lang,
			Locale: session.Locale},
		UserAgent:   ctx.UserAgent,
		LoginIP:     ctx.IP,
		LoginRegion: ctx.Region}

	_, err := sb.WriteModel.Client.HandleLogin(context.Background(), &req)
	if err != nil {
		log.Errorf("Register login error, %v", err)
	}
	return err
}
