package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gavrilaf/spawn/pkg/env"
	mdl "github.com/gavrilaf/spawn/pkg/model"
)

type RedisCache struct {
	conn redis.Conn
}

func Connect(en *env.Environment) (*RedisCache, error) {

	//redis://[:password@]host[:port][/db-number][?option=value]
	conn, err := redis.DialURL("redis://localhost:7001")
	if err != nil {
		return nil, err
	}
	return &RedisCache{conn}, nil
}

func (cache *RedisCache) Close() error {
	return cache.conn.Close()
}

type UserCache interface {
	AddClient(client mdl.Client) error
	FindClient(id string) (*mdl.Client, error)

	AddSession(session mdl.Session) error
	FindSession(id string) (*mdl.Session, error)
	DeleteSession(id string) error

	AddUser(profile mdl.UserProfile, devices []string) error
	FindProfile(id string) (*mdl.UserProfile, error)

	AddDevice(userId string, deviceId string) error
	DeleteDevice(userId string, deviceId string) error
	IsDeviceExists(userId string, deviceId string) (bool, error)
}
