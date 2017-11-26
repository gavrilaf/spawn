package auth

import (
	"github.com/gavrilaf/spawn/pkg/api"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/model"
)

type Storage interface {
	Close()

	FindClient(id string) (*db.Client, error)

	RegisterUser(username string, password string, device db.DeviceInfo) error

	FindUser(username string) (*mdl.AuthUser, error)

	FindDevice(userId string, deviceId string) (*mdl.AuthDevice, error)
	AddDevice(userId string, device db.DeviceInfo) (*mdl.AuthDevice, error)

	StoreSession(session mdl.Session) error
	FindSession(id string) (*mdl.Session, error)

	HandlerLogin(session mdl.Session, ctx LoginContext) error
}

type StorageImpl struct {
	*api.Bridge
}
