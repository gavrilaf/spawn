package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/env"
)

const (
	Scope = "read-model"
)

// TODO: Rename package to rmd (read model)

// Connect to the spawn read model
func Connect(en *env.Environment) (Cache, error) {

	//redis://[:password@]host[:port][/db-number][?option=value]
	conn, err := redis.DialURL("redis://localhost:7001")
	if err != nil {
		return nil, err
	}
	return &Bridge{conn}, nil
}

type Cache interface {
	Close() error

	// Auth cache

	AddClient(client db.Client) error
	FindClient(id string) (*db.Client, error)

	AddSession(session mdl.Session) error
	FindSession(id string) (*mdl.Session, error)
	DeleteSession(id string) error

	SetUserAuthInfo(profile db.UserProfile, devices []db.DeviceInfo) error
	FindUserAuthInfo(username string) (*mdl.AuthUser, error)

	SetDevice(device db.DeviceInfo) error
	DeleteDevice(userId string, deviceId string) error
	FindDevice(userId string, deviceId string) (*mdl.AuthDevice, error)

	AddConfirmCode(kind string, id string, code string) error
	GetConfirmCode(kind string, id string) (string, error)
	DeleteConfirmCode(kind string, id string) error

	// User profile

	SetUserProfile(profile db.UserProfile) error
	GetUserProfile(userID string) (*mdl.UserProfile, error)

	SetUserDevicesInfo(userID string, devices []db.DeviceInfoEx) error
	GetUserDevicesInfo(userID string) ([]mdl.UserDeviceInfo, error)
}

//////////////////////////////////////////////////////////////////////////////////////////

type Bridge struct {
	conn redis.Conn
}

func (cache *Bridge) Close() error {
	if cache == nil || cache.conn == nil {
		return nil
	}

	return cache.conn.Close()
}

//////////////////////////////////////////////////////////////////////////////////////////

func (p *Bridge) getKeys(pattern string) ([]string, error) {
	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(p.conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}
