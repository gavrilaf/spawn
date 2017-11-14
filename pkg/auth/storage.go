package auth

import (
	"github.com/gavrilaf/spawn/pkg/cache"
	mdl "github.com/gavrilaf/spawn/pkg/model"
)

type Storage interface {
	FindClient(id string) (mdl.Client, error)

	RegisterUser(username string, password string, device mdl.DeviceInfo) error

	FindUser(username string) (*cache.AuthUser, error)

	FindDevice(userId string, deviceId string) (*cache.AuthDevice, error)
	AddDevice(userId string, device mdl.DeviceInfo) (*cache.AuthDevice, error)

	StoreSession(session cache.Session) error
	FindSession(id string) (*cache.Session, error)
}
