package cache

import (
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func userDevicesID(userID string, deviceID string) string {
	return userDevicesPattern(userID) + deviceID
}

func userDevicesPattern(userID string) string {
	return "userdex:" + userID + "-"
}

func (br *Bridge) SetUserDevicesInfo(userID string, devices []db.DeviceInfoEx) error {
	conn := br.get()
	defer conn.Close()

	if err := br.deleteUserDevicesInfo(userID); err != nil {
		return err
	}

	for _, d := range devices {
		key := userDevicesID(userID, d.ID)
		dd := mdl.CreateUserDeviceInfoFromDb(d)
		_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&dd)...)

		if err != nil {
			br.deleteUserDevicesInfo(userID)
			return err
		}

	}
	return nil
}

func (br *Bridge) GetUserDevicesInfo(userID string) ([]mdl.UserDeviceInfo, error) {
	conn := br.get()
	defer conn.Close()

	keys, err := br.getKeys(userDevicesPattern(userID) + "*")
	if err != nil {
		return nil, err
	}

	devices := make([]mdl.UserDeviceInfo, len(keys))

	for indx, key := range keys {
		v, err := redis.Values(conn.Do("HGETALL", key))

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

func (br *Bridge) deleteUserDevicesInfo(userID string) error {
	conn := br.get()
	defer conn.Close()

	keys, err := br.getKeys(userDevicesPattern(userID) + "*")
	if err != nil {
		return err
	}

	for _, key := range keys {
		conn.Do("DEL", key)
	}

	return nil
}
