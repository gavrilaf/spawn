package cache

import (
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func profileID(id string) string {
	return "profile:" + id
}

func (cache *Bridge) SetUserProfile(profile db.UserProfile) error {
	cacheProfile := mdl.CreateProfileFromDbModel(profile)
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(profileID(cacheProfile.ID)).AddFlat(&cacheProfile)...)
	return err
}

func (cache *Bridge) GetUserProfile(userID string) (*mdl.UserProfile, error) {
	key := profileID(userID)
	v, err := redis.Values(cache.conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, key)
	}

	var profile mdl.UserProfile
	if err := redis.ScanStruct(v, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func userDevicesID(userID string, deviceID string) string {
	return userDevicesPattern(userID) + deviceID
}

func userDevicesPattern(userID string) string {
	return "userdex:" + userID + "-"
}

func (cache *Bridge) SetUserDevicesInfo(userID string, devices []db.DeviceInfoEx) error {
	if err := cache.deleteUserDevicesInfo(userID); err != nil {
		return err
	}

	for _, d := range devices {
		key := userDevicesID(userID, d.ID)
		dd := mdl.CreateUserDeviceInfoFromDb(d)
		_, err := cache.conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&dd)...)

		if err != nil {
			cache.deleteUserDevicesInfo(userID)
			return err
		}

	}
	return nil
}

func (cache *Bridge) GetUserDevicesInfo(userID string) ([]mdl.UserDeviceInfo, error) {
	keys, err := cache.getKeys(userDevicesPattern(userID) + "*")
	if err != nil {
		return nil, err
	}

	devices := make([]mdl.UserDeviceInfo, len(keys))

	for indx, key := range keys {
		v, err := redis.Values(cache.conn.Do("HGETALL", key))

		if err != nil {
			return nil, err
		}
		if len(v) == 0 {
			return nil, errx.ErrKeyNotFound(Scope, key)
		}

		var device mdl.UserDeviceInfo
		if err := redis.ScanStruct(v, &device); err != nil {
			return nil, err
		}
		devices[indx] = device
	}

	return devices, nil
}

func (cache *Bridge) deleteUserDevicesInfo(userID string) error {
	keys, err := cache.getKeys(userDevicesPattern(userID) + "*")
	if err != nil {
		return err
	}

	for _, key := range keys {
		cache.conn.Do("DEL", key)
	}

	return nil
}
