package auth

import (
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
)

func generateRefreshToken() string {
	k, _ := cryptx.GenerateSaltedKey(uuid.NewV4().String())
	return hex.EncodeToString(k)
}

func createLoginContext(c *gin.Context) LoginContext {
	return LoginContext{
		IP:        c.ClientIP(),
		UserAgent: c.Request.Header.Get("User-Agent"),
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////

func (p ApiImpl) getClient(id string) (*db.Client, error) {
	return p.ReadModel.FindClient(id)
}

func (p ApiImpl) registerUser(username string, password string, device db.DeviceInfo) error {
	req := pb.CreateUserReq{
		Username:     username,
		PasswordHash: password,
		Device: &pb.Device{
			ID:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	resp, err := p.WriteModel.CreateUser(&req)
	if err == nil && resp != nil {
		log.Infof("Registered user (%s, %s)", username, resp.ID)
	}

	return err
}

func (self ApiImpl) findUser(username string) (*mdl.AuthUser, error) {
	return self.ReadModel.FindUserAuthInfo(username)
}

func (self ApiImpl) getDevice(userId string, deviceId string) (*mdl.AuthDevice, error) {
	return self.ReadModel.GetDevice(userId, deviceId)
}

func (self ApiImpl) addDevice(userID string, device db.DeviceInfo) (*mdl.AuthDevice, error) {
	req := pb.UserDevice{
		UserID: userID,
		Device: &pb.Device{
			ID:     device.ID,
			Name:   device.Name,
			Locale: device.Locale,
			Lang:   device.Lang},
	}

	_, err := self.WriteModel.AddDevice(&req)
	if err == nil {
		log.Infof("Added device (%s, %s)", userID, device.ID)
	}

	return self.getDevice(userID, device.ID)
}

func (self ApiImpl) addSession(session mdl.Session) (string, error) {
	return self.ReadModel.AddSession(session, false)
}

func (self ApiImpl) updateSession(session mdl.Session) error {
	return self.ReadModel.SetSession(session)
}

func (self ApiImpl) getSession(id string) (*mdl.Session, error) {
	session, err := self.ReadModel.GetSession(id)
	if err != nil {
		log.Errorf("Can't find session with id %s, error: %v", id, err)
		return nil, types.ErrSessionNotFound
	}

	return session, nil
}

func (self ApiImpl) handlerLogin(session mdl.Session, ctx LoginContext) error {
	req := pb.LoginReq{
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

	_, err := self.WriteModel.HandleLogin(&req)
	if err != nil {
		log.Errorf("Register login error, %v", err)
	}
	return err
}
