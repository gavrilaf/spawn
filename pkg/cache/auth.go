package cache

import (
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	confirmExpiration = 30 * 60 // 30 minutes
)

// Client

func clientRedisID(id string) string {
	return "client:" + id
}

func (cache *Bridge) AddClient(client db.Client) error {
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(clientRedisID(client.ID)).AddFlat(&client)...)
	return err
}

func (cache *Bridge) FindClient(id string) (*db.Client, error) {
	key := clientRedisID(id)
	v, err := redis.Values(cache.conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, key)
	}

	var client db.Client
	if err := redis.ScanStruct(v, &client); err != nil {
		return nil, err
	}

	return &client, nil
}

// Session

func sessionRedisID(id string) string {
	return "session:" + id
}

func (cache *Bridge) AddSession(session mdl.Session) error {
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(sessionRedisID(session.ID)).AddFlat(&session)...)
	return err
}

func (cache *Bridge) FindSession(id string) (*mdl.Session, error) {
	key := sessionRedisID(id)
	v, err := redis.Values(cache.conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, key)
	}

	var session mdl.Session
	if err := redis.ScanStruct(v, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (cache *Bridge) DeleteSession(id string) error {
	_, err := cache.conn.Do("DEL", sessionRedisID(id))
	return err
}

// Users
func authUserID(username string) string {
	return "user:" + username
}

func (cache *Bridge) SetUserAuthInfo(profile db.UserProfile, devices []db.DeviceInfo) error {
	authUser := mdl.CreateAuthUserFromProfile(profile)

	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(authUserID(profile.Username)).AddFlat(&authUser)...)
	if err != nil {
		return err
	}

	for _, d := range devices {
		d.UserID = profile.ID
		err = cache.SetDevice(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cache *Bridge) FindUserAuthInfo(username string) (*mdl.AuthUser, error) {
	key := authUserID(username)
	v, err := redis.Values(cache.conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, key)
	}

	var user mdl.AuthUser
	if err := redis.ScanStruct(v, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Devices

func authDeviceID(userID string, deviceID string) string {
	return "device:" + userID + deviceID
}

func (cache *Bridge) SetDevice(device db.DeviceInfo) error {
	ad := mdl.CreateAuthDeviceFromDevice(device)
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(authDeviceID(device.UserID, device.ID)).AddFlat(&ad)...)
	return err
}

func (cache *Bridge) DeleteDevice(userID string, deviceID string) error {
	_, err := cache.conn.Do("DEL", authDeviceID(userID, deviceID))
	return err
}

func (cache *Bridge) FindDevice(userID string, deviceID string) (*mdl.AuthDevice, error) {
	key := authDeviceID(userID, deviceID)
	v, err := redis.Values(cache.conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, key)
	}

	var d mdl.AuthDevice
	if err := redis.ScanStruct(v, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

// Confirm code
func (cache *Bridge) AddConfirmCode(kind string, id string, code string) error {
	key := "confirm:" + kind + id
	_, err := cache.conn.Do("SETEX", key, confirmExpiration, code)
	return err
}

func (cache *Bridge) GetConfirmCode(kind string, id string) (string, error) {
	key := "confirm:" + kind + id
	return redis.String(cache.conn.Do("GET", key))
}

func (cache *Bridge) DeleteConfirmCode(kind string, id string) error {
	key := "confirm:" + kind + id
	_, err := cache.conn.Do("DEL", key)
	return err
}
