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

func (cache *RedisCache) AddClient(client mdl.Client) error {
	_, err := cache.conn.Do("SET", clientRedisId(client.ID), client.Secret)
	return err
}

func (cache *RedisCache) FindClient(id string) (*mdl.Client, error) {
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

func (cache *RedisCache) AddSession(session mdl.Session) error {
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(sessionRedisId(session.ID)).AddFlat(&session)...)
	return err
}

func (cache *RedisCache) FindSession(id string) (*mdl.Session, error) {
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

func (cache *RedisCache) DeleteSession(id string) error {
	_, err := cache.conn.Do("DEL", sessionRedisId(id))
	return err
}

// Users
func profileRedisId(id string) string {
	return "profile:" + id
}

func devicesRedisId(id string) string {
	return "devices:" + id
}

func (cache *RedisCache) AddUser(profile mdl.UserProfile, devices []string) error {
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(profileRedisId(profile.ID)).AddFlat(&profile)...)
	if err != nil {
		return err
	}

	_, err = cache.conn.Do("SADD", redis.Args{}.Add(devicesRedisId(profile.ID)).AddFlat(devices)...)
	return err

}

func (cache *RedisCache) FindProfile(id string) (*mdl.UserProfile, error) {
	v, err := redis.Values(cache.conn.Do("HGETALL", profileRedisId(id)))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errNotFound(profileRedisId(id))
	}

	var profile mdl.UserProfile
	if err := redis.ScanStruct(v, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// Devices
func (cache *RedisCache) AddDevice(userId string, deviceId string) error {
	_, err := cache.conn.Do("SADD", redis.Args{}.Add(devicesRedisId(userId)).Add(deviceId)...)
	return err
}

func (cache *RedisCache) DeleteDevice(userId string, deviceId string) error {
	_, err := cache.conn.Do("SREM", redis.Args{}.Add(devicesRedisId(userId)).Add(deviceId)...)
	return err
}

func (cache *RedisCache) IsDeviceExists(userId string, deviceId string) (bool, error) {
	return redis.Bool(cache.conn.Do("SISMEMBER", redis.Args{}.Add(devicesRedisId(userId)).Add(deviceId)...))
}
