package cache

import (
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	Scope                  = "read-model"
	ReasonSessionDuplicate = "session-already-exist"
)

var ErrSessionDuplicate = errx.New(Scope, ReasonSessionDuplicate)

// Connect to the spawn read model
func Connect(en *env.Environment) (Cache, error) {

	//redis://[:password@]host[:port][/db-number][?option=value]
	return &Bridge{newPool(en)}, nil
}

type Cache interface {
	Close() error

	// Auth

	AddClient(client db.Client) error
	FindClient(id string) (*db.Client, error)

	// Save Session to the Read model, check if session already exists for pair (user, device).
	// Only one session allowed for (user, device) pair
	// forced - if session exists, remove it and create new
	// returns stored session id
	AddSession(session mdl.Session, forced bool) (string, error)

	SetSession(session mdl.Session) error
	GetSession(id string) (*mdl.Session, error)
	DeleteSession(id string) error

	// find session for (user, device) pair
	FindSession(userID string, deviceID string) (*mdl.Session, error)

	SetUserAuthInfo(profile db.UserProfile, devices []db.DeviceInfo) error
	FindUserAuthInfo(username string) (*mdl.AuthUser, error)

	SetDevice(device db.DeviceInfo) error
	DeleteDevice(userID string, deviceID string) error
	GetDevice(userID string, deviceID string) (*mdl.AuthDevice, error)

	// User

	SetUserDevicesInfo(userID string, devices []db.DeviceInfoEx) error
	GetUserDevicesInfo(userID string) ([]mdl.UserDeviceInfo, error)

	AddConfirmCode(kind string, id string, code string) error
	GetConfirmCode(kind string, id string) (string, error)
	DeleteConfirmCode(kind string, id string) error

	// User profile

	SetUserProfile(profile db.UserProfile) error
	GetUserProfile(userID string) (*mdl.UserProfile, error)
}
