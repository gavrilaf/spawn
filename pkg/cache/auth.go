package cache

import (
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/model"
)

const (
	confirmExpiration = 30 * 60 // 30 minutes
)

// Client

func clientRedisId(id string) string {
	return "client:" + id
}

func (cache *Bridge) AddClient(client db.Client) error {
	_, err := cache.conn.Do("SET", clientRedisId(client.ID), client.Secret)
	return err
}

func (cache *Bridge) FindClient(id string) (*db.Client, error) {
	secret, err := redis.Bytes(cache.conn.Do("GET", clientRedisId(id)))
	if err != nil {
		return nil, err
	}
	return &db.Client{id, secret}, nil
}

// Session

func sessionRedisId(id string) string {
	return "session:" + id
}

func (cache *Bridge) AddSession(session mdl.Session) error {
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(sessionRedisId(session.ID)).AddFlat(&session)...)
	return err
}

func (cache *Bridge) FindSession(id string) (*mdl.Session, error) {
	v, err := redis.Values(cache.conn.Do("HGETALL", sessionRedisId(id)))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errNotFound(sessionRedisId(id))
	}

	var session mdl.Session
	if err := redis.ScanStruct(v, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (cache *Bridge) DeleteSession(id string) error {
	_, err := cache.conn.Do("DEL", sessionRedisId(id))
	return err
}

// Users
func authUserId(username string) string {
	return "user:" + username
}

func (cache *Bridge) SetUserAuthInfo(profile db.UserProfile, devices []db.DeviceInfo) error {
	authUser := mdl.CreateAuthUserFromProfile(profile)

	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(authUserId(profile.Username)).AddFlat(&authUser)...)
	if err != nil {
		return err
	}

	for _, d := range devices {
		err = cache.SetDevice(profile.ID, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cache *Bridge) FindUserAuthInfo(username string) (*mdl.AuthUser, error) {
	v, err := redis.Values(cache.conn.Do("HGETALL", authUserId(username)))

	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	var user mdl.AuthUser
	if err := redis.ScanStruct(v, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Devices

func authDeviceId(userId string, deviceId string) string {
	return "device:" + userId + deviceId
}

func (cache *Bridge) SetDevice(userID string, device db.DeviceInfo) error {
	device.UserID = userID
	ad := mdl.CreateAuthDeviceFromDevice(device)
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(authDeviceId(userID, ad.DeviceID)).AddFlat(&ad)...)
	return err
}

func (cache *Bridge) DeleteDevice(userId string, deviceId string) error {
	_, err := cache.conn.Do("DEL", authDeviceId(userId, deviceId))
	return err
}

func (cache *Bridge) FindDevice(userId string, deviceId string) (*mdl.AuthDevice, error) {
	v, err := redis.Values(cache.conn.Do("HGETALL", authDeviceId(userId, deviceId)))
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
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
