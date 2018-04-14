package cache

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
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

func sessionPatern(userID string, deviceID string) string {
	return "session:" + userID + ":" + deviceID
}

func (br *Bridge) AddSession(session mdl.Session, forced bool) (string, error) {
	conn := br.get()
	defer conn.Close()

	// check old session
	keys, err := getKeys(conn, sessionPatern(session.UserID, session.DeviceID)+"*")
	if err != nil {
		return "", err
	}

	if len(keys) > 0 {
		if !forced {
			return "", ErrSessionDuplicate
		}

		// invalidate old session
		for _, key := range keys {
			conn.Do("DEL", key)
		}
	}

	nonce, _ := redis.Int(conn.Do("INCR", "sessions-counter"))
	key := sessionPatern(session.UserID, session.DeviceID) + ":" + strconv.Itoa(nonce)

	session.ID = key
	_, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&session)...)
	if err != nil {
		return "", nil
	}

	return key, err
}

func (br *Bridge) SetSession(session mdl.Session) error {
	conn := br.get()
	defer conn.Close()

	_, err := conn.Do("HMSET", redis.Args{}.Add(session.ID).AddFlat(&session)...)
	return err
}

func (br *Bridge) GetSession(id string) (*mdl.Session, error) {
	conn := br.get()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETALL", id))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, id)
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

	_, err := conn.Do("DEL", id)
	return err
}

func (br *Bridge) FindSession(userID string, deviceID string) (*mdl.Session, error) {
	conn := br.get()
	defer conn.Close()

	key := sessionPatern(userID, deviceID)
	keys, err := getKeys(conn, key+"*")
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, key)
	}

	key = keys[0]
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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func authUserID(username string) string {
	return "auser:" + username
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
	return "adevice:" + userID + deviceID
}

func (br *Bridge) SetDevice(device db.DeviceInfo) error {
	conn := br.get()
	defer conn.Close()

	ad := mdl.CreateAuthDeviceFromDevice(device)
	_, err := conn.Do("HMSET", redis.Args{}.Add(authDeviceID(device.UserID, device.DeviceID)).AddFlat(&ad)...)
	return err
}

func (br *Bridge) DeleteDevice(userID string, deviceID string) error {
	conn := br.get()
	defer conn.Close()

	_, err := conn.Do("DEL", authDeviceID(userID, deviceID))
	return err
}

func (br *Bridge) GetDevice(userID string, deviceID string) (*mdl.AuthDevice, error) {
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
