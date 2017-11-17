package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gavrilaf/spawn/pkg/env"
	mdl "github.com/gavrilaf/spawn/pkg/model"
)

type Cache struct {
	conn redis.Conn
}

func Connect(en *env.Environment) (*Cache, error) {

	//redis://[:password@]host[:port][/db-number][?option=value]
	conn, err := redis.DialURL("redis://localhost:7001")
	if err != nil {
		return nil, err
	}
	return &Cache{conn}, nil
}

func (cache *Cache) Close() error {
	return cache.conn.Close()
}

type UserCache interface {
	AddClient(client mdl.Client) error
	FindClient(id string) (*mdl.Client, error)

	AddSession(session Session) error
	FindSession(id string) (*Session, error)
	DeleteSession(id string) error

	SetUserAuthInfo(profile mdl.UserProfile, devices []mdl.DeviceInfo) error
	FindUserAuthInfo(username string) (*AuthUser, error)

	SetDevice(userID string, device mdl.DeviceInfo) error
	DeleteDevice(userId string, deviceId string) error
	FindDevice(userId string, deviceId string) (*AuthDevice, error)

	AddConfirmCode(kind string, id string, code string) error
	GetConfirmCode(kind string, id string) (string, error)
	DeleteConfirmCode(kind string, id string) error
}
