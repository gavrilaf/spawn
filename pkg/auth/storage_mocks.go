package auth

import (
	"fmt"
	"github.com/gavrilaf/spawn/pkg/cache"
	"github.com/gavrilaf/spawn/pkg/env"
	mdl "github.com/gavrilaf/spawn/pkg/model"
)

/*
 * Clients storage
 */

var DefaultClients = map[string][]byte{
	"client_test": []byte("client_test_key"),
	"client_ios":  []byte("client_ios_key"),
}

var DefaultUsers = []cache.AuthUser{
	cache.AuthUser{
		ID: "user1",
		AuthInfo: mdl.AuthInfo{
			Username:     "id1@spawn.com",
			PasswordHash: "24326124313024636d636548593741794976623167777154732f502f754c4c466a4c755543784a6b696f386f6b4c344a42686a514b76494943654653",
			Permissions: mdl.Permissions{
				IsLocked:         false,
				IsEmailConfirmed: false,
				Is2FARequired:    false}}},

	cache.AuthUser{
		ID: "user2",
		AuthInfo: mdl.AuthInfo{
			Username:     "id2@spawn.com",
			PasswordHash: "243261243130247a6b4576684654664b4c4a5945486871766e6b51472e354771335676664a492f6a5232304c73465774354c553939314b3944766b6d",
			Permissions: mdl.Permissions{
				IsLocked:         false,
				IsEmailConfirmed: false,
				Is2FARequired:    false}}},
}

var DefaultDevices = []cache.AuthDevice{
	cache.AuthDevice{DeviceID: "d1", UserID: "user1", IsConfirmed: true, Locale: "en", Lang: "en"},
	cache.AuthDevice{DeviceID: "d2", UserID: "user1", IsConfirmed: true, Locale: "en", Lang: "en"},
	cache.AuthDevice{DeviceID: "d2", UserID: "user2", IsConfirmed: true, Locale: "en", Lang: "en"},
	cache.AuthDevice{DeviceID: "d3", UserID: "user2", IsConfirmed: false, Locale: "ru", Lang: "ru"},
}

type StorageMock struct {
	counter  int
	users    map[string]cache.AuthUser
	devices  map[string]map[string]cache.AuthDevice
	sessions map[string]cache.Session
}

func NewStorageMock(en *env.Environment) *StorageMock {
	users := make(map[string]cache.AuthUser)
	devices := make(map[string]map[string]cache.AuthDevice)

	for _, v := range DefaultUsers {
		users[v.Username] = v
	}

	for _, d := range DefaultDevices {
		devices[d.UserID] = make(map[string]cache.AuthDevice)
		devices[d.UserID][d.DeviceID] = d
	}

	return &StorageMock{counter: 3, users: users, devices: devices, sessions: make(map[string]cache.Session)}
}

///////

func (c *StorageMock) FindClient(id string) (mdl.Client, error) {
	secret, ok := DefaultClients[id]
	if !ok {
		return mdl.Client{}, errClientUnknown
	}
	return mdl.Client{ID: id, Secret: secret}, nil
}

func (c *StorageMock) RegisterUser(username string, password string, device mdl.DeviceInfo) error {
	if user, _ := c.FindUser(username); user != nil {
		return errUserAlreadyExist
	}

	id := fmt.Sprintf("id%d", c.counter)
	c.counter += 1
	user := cache.AuthUser{
		ID: id,
		AuthInfo: mdl.AuthInfo{
			Username:     username,
			PasswordHash: password,
			Permissions: mdl.Permissions{
				IsLocked:         false,
				IsEmailConfirmed: false,
				Is2FARequired:    false}}}

	c.users[username] = user

	nd := cache.CreateAuthDeviceFromDevice(device)
	nd.IsConfirmed = true
	c.devices[user.ID] = make(map[string]cache.AuthDevice)
	c.devices[user.ID][device.ID] = nd

	return nil
}

func (c *StorageMock) AddDevice(userId string, device mdl.DeviceInfo) (*cache.AuthDevice, error) {
	d := cache.CreateAuthDeviceFromDevice(device)
	c.devices[userId][device.ID] = d
	return &d, nil
}

func (c *StorageMock) FindUser(username string) (*cache.AuthUser, error) {
	user, ok := c.users[username]
	if !ok {
		return nil, nil
	}

	return &user, nil
}

func (c *StorageMock) FindDevice(userId string, deviceId string) (*cache.AuthDevice, error) {
	devices, ok := c.devices[userId]
	if !ok {
		return nil, nil
	}

	d, ok := devices[deviceId]
	if !ok {
		return nil, nil
	}

	return &d, nil
}

func (c *StorageMock) StoreSession(session cache.Session) error {
	c.sessions[session.ID] = session
	return nil
}

func (c *StorageMock) FindSession(id string) (*cache.Session, error) {
	session, ok := c.sessions[id]
	if !ok {
		return nil, errSessionNotFound
	}
	return &session, nil
}
