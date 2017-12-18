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

func (br *Bridge) AddClient(client db.Client) error {
	conn := br.get()
	defer conn.Close()

	_, err := conn.Do("HMSET", redis.Args{}.Add(clientRedisID(client.ID)).AddFlat(&client)...)
	return err
}

func (br *Bridge) FindClient(id string) (*db.Client, error) {
	conn := br.get()
	defer conn.Close()

	key := clientRedisID(id)
	v, err := redis.Values(conn.Do("HGETALL", key))

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

func (br *Bridge) AddSession(session mdl.Session) error {
	conn := br.get()
	defer conn.Close()

	_, err := conn.Do("HMSET", redis.Args{}.Add(sessionRedisID(session.ID)).AddFlat(&session)...)
	return err
}

func (br *Bridge) FindSession(id string) (*mdl.Session, error) {
	conn := br.get()
	defer conn.Close()

	key := sessionRedisID(id)
	v, err := redis.Values(conn.Do("HGETALL", key))

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

func (br *Bridge) DeleteSession(id string) error {
	conn := br.get()
	defer conn.Close()

	_, err := conn.Do("DEL", sessionRedisID(id))
	return err
}

// Users
func authUserID(username string) string {
	return "user:" + username
}

func (br *Bridge) SetUserAuthInfo(profile db.UserProfile, devices []db.DeviceInfo) error {
	conn := br.get()
	defer conn.Close()

	authUser := mdl.CreateAuthUserFromProfile(profile)

	_, err := conn.Do("HMSET", redis.Args{}.Add(authUserID(profile.Username)).AddFlat(&authUser)...)
	if err != nil {
		return err
	}

	for _, d := range devices {
		d.UserID = profile.ID
		err = br.SetDevice(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (br *Bridge) FindUserAuthInfo(username string) (*mdl.AuthUser, error) {
	conn := br.get()
	defer conn.Close()

	key := authUserID(username)
	v, err := redis.Values(conn.Do("HGETALL", key))

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

func (br *Bridge) SetDevice(device db.DeviceInfo) error {
	conn := br.get()
	defer conn.Close()

	ad := mdl.CreateAuthDeviceFromDevice(device)
	_, err := conn.Do("HMSET", redis.Args{}.Add(authDeviceID(device.UserID, device.ID)).AddFlat(&ad)...)
	return err
}

func (br *Bridge) DeleteDevice(userID string, deviceID string) error {
	conn := br.get()
	defer conn.Close()

	_, err := conn.Do("DEL", authDeviceID(userID, deviceID))
	return err
}

func (br *Bridge) FindDevice(userID string, deviceID string) (*mdl.AuthDevice, error) {
	conn := br.get()
	defer conn.Close()

	key := authDeviceID(userID, deviceID)
	v, err := redis.Values(conn.Do("HGETALL", key))

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
func (br *Bridge) AddConfirmCode(kind string, id string, code string) error {
	conn := br.get()
	defer conn.Close()

	key := "confirm:" + kind + id
	_, err := conn.Do("SETEX", key, confirmExpiration, code)
	return err
}

func (br *Bridge) GetConfirmCode(kind string, id string) (string, error) {
	conn := br.get()
	defer conn.Close()

	key := "confirm:" + kind + id
	return redis.String(conn.Do("GET", key))
}

func (br *Bridge) DeleteConfirmCode(kind string, id string) error {
	conn := br.get()
	defer conn.Close()

	key := "confirm:" + kind + id
	_, err := conn.Do("DEL", key)
	return err
}
