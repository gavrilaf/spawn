package auth

import (
	mdl "github.com/gavrilaf/spawn/pkg/model"
)

type Storage interface {
	FindClient(id string) (*mdl.Client, error)

	RegisterUser(username string, password string, deviceId string) error

	FindUser(username string) (*mdl.UserProfile, error)
	IsDeviceAllowed(userId string, deviceId string) (bool, error)

	StoreSession(session mdl.Session) error
	FindSession(id string) (*mdl.Session, error)
}
