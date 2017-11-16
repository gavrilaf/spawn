package cache

import (
	//"fmt"
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/model"
)

// Client

func clientRedisId(id string) string {
	return "client:" + id
}

func (cache *Cache) AddClient(client mdl.Client) error {
	_, err := cache.conn.Do("SET", clientRedisId(client.ID), client.Secret)
	return err
}

func (cache *Cache) FindClient(id string) (*mdl.Client, error) {
	secret, err := redis.Bytes(cache.conn.Do("GET", clientRedisId(id)))
	if err != nil {
		return nil, err
	}
	return &mdl.Client{id, secret}, nil
}

// Session

func sessionRedisId(id string) string {
	return "session:" + id
}

func (cache *Cache) AddSession(session Session) error {
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(sessionRedisId(session.ID)).AddFlat(&session)...)
	return err
}

func (cache *Cache) FindSession(id string) (*Session, error) {
	v, err := redis.Values(cache.conn.Do("HGETALL", sessionRedisId(id)))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errNotFound(sessionRedisId(id))
	}

	var session Session
	if err := redis.ScanStruct(v, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (cache *Cache) DeleteSession(id string) error {
	_, err := cache.conn.Do("DEL", sessionRedisId(id))
	return err
}

// Users
func authUserId(username string) string {
	return "user:" + username
}

func authDeviceId(userId string, deviceId string) string {
	return "device:" + userId + deviceId
}

func (cache *Cache) AddUserAuthInfo(profile mdl.UserProfile, devices []mdl.DeviceInfo) error {
	authUser := CreateAuthUserFromProfile(profile)

	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(authUserId(profile.Username)).AddFlat(&authUser)...)
	if err != nil {
		return err
	}

	for _, d := range devices {
		err = cache.AddDevice(profile.ID, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cache *Cache) FindUserAuthInfo(username string) (*AuthUser, error) {
	v, err := redis.Values(cache.conn.Do("HGETALL", authUserId(username)))

	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	var user AuthUser
	if err := redis.ScanStruct(v, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Devices

func (cache *Cache) AddDevice(userID string, device mdl.DeviceInfo) error {
	device.UserID = userID
	ad := CreateAuthDeviceFromDevice(device)
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(authDeviceId(userID, ad.DeviceID)).AddFlat(&ad)...)
	return err
}

func (cache *Cache) DeleteDevice(userId string, deviceId string) error {
	_, err := cache.conn.Do("DEL", authDeviceId(userId, deviceId))
	return err
}

func (cache *Cache) FindDevice(userId string, deviceId string) (*AuthDevice, error) {
	v, err := redis.Values(cache.conn.Do("HGETALL", authDeviceId(userId, deviceId)))
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	var d AuthDevice
	if err := redis.ScanStruct(v, &d); err != nil {
		return nil, err
	}

	return &d, nil
}
