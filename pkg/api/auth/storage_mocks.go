package auth

import (
	"fmt"

	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/env"
)

var storageMock = NewStorageMock(env.GetEnvironment("Test"))

/*
 * Clients storage
 */

var DefaultClients = map[string][]byte{
	"client_test": []byte("client_test_key"),
	"client_ios":  []byte("client_ios_key"),
}

var DefaultUsers = []mdl.AuthUser{
	mdl.AuthUser{
		ID: "user1",
		AuthInfo: db.AuthInfo{
			Username:     "id1@spawn.com",
			PasswordHash: "24326124313024636d636548593741794976623167777154732f502f754c4c466a4c755543784a6b696f386f6b4c344a42686a514b76494943654653",
			Permissions: db.Permissions{
				IsLocked:         false,
				IsEmailConfirmed: false,
				Is2FARequired:    false}}},

	mdl.AuthUser{
		ID: "user2",
		AuthInfo: db.AuthInfo{
			Username:     "id2@spawn.com",
			PasswordHash: "243261243130247a6b4576684654664b4c4a5945486871766e6b51472e354771335676664a492f6a5232304c73465774354c553939314b3944766b6d",
			Permissions: db.Permissions{
				IsLocked:         false,
				IsEmailConfirmed: false,
				Is2FARequired:    false}}},
}

var DefaultDevices = []mdl.AuthDevice{
	mdl.AuthDevice{DeviceID: "d1", UserID: "user1", IsConfirmed: true, Locale: "en", Lang: "en"},
	mdl.AuthDevice{DeviceID: "d2", UserID: "user1", IsConfirmed: true, Locale: "en", Lang: "en"},
	mdl.AuthDevice{DeviceID: "d2", UserID: "user2", IsConfirmed: true, Locale: "en", Lang: "en"},
	mdl.AuthDevice{DeviceID: "d3", UserID: "user2", IsConfirmed: false, Locale: "ru", Lang: "ru"},
}

type StorageMock struct {
	counter  int
	users    map[string]mdl.AuthUser
	devices  map[string]map[string]mdl.AuthDevice
	sessions map[string]mdl.Session
}

func NewStorageMock(en *env.Environment) *StorageMock {
	users := make(map[string]mdl.AuthUser)
	devices := make(map[string]map[string]mdl.AuthDevice)

	for _, v := range DefaultUsers {
		users[v.Username] = v
		devices[v.ID] = make(map[string]mdl.AuthDevice)
	}

	for _, d := range DefaultDevices {
		devices[d.UserID][d.DeviceID] = d
	}

	return &StorageMock{counter: 3, users: users, devices: devices, sessions: make(map[string]mdl.Session)}
}

///////
func (c *StorageMock) Close() {}

func (c *StorageMock) FindClient(id string) (*db.Client, error) {
	secret, ok := DefaultClients[id]
	if !ok {
		return nil, errClientUnknown
	}
	return &db.Client{ID: id, Secret: secret}, nil
}

func (c *StorageMock) RegisterUser(username string, password string, device db.DeviceInfo) error {
	if user, _ := c.FindUser(username); user != nil {
		return errUserAlreadyExist
	}

	id := fmt.Sprintf("id%d", c.counter)
	c.counter += 1
	user := mdl.AuthUser{
		ID: id,
		AuthInfo: db.AuthInfo{
			Username:     username,
			PasswordHash: password,
			Permissions: db.Permissions{
				IsLocked:         false,
				IsEmailConfirmed: false,
				Is2FARequired:    false}}}

	c.users[username] = user

	nd := mdl.CreateAuthDeviceFromDevice(device)
	nd.IsConfirmed = true
	c.devices[user.ID] = make(map[string]mdl.AuthDevice)
	c.devices[user.ID][device.ID] = nd

	return nil
}

func (c *StorageMock) AddDevice(userId string, device db.DeviceInfo) (*mdl.AuthDevice, error) {
	d := mdl.CreateAuthDeviceFromDevice(device)
	c.devices[userId][device.ID] = d
	return &d, nil
}

func (c *StorageMock) FindUser(username string) (*mdl.AuthUser, error) {
	user, ok := c.users[username]
	if !ok {
		return nil, errUserUnknown
	}

	return &user, nil
}

func (c *StorageMock) FindDevice(userId string, deviceId string) (*mdl.AuthDevice, error) {
	devices, ok := c.devices[userId]
	if !ok {
		return nil, errDeviceUnknown
	}

	d, ok := devices[deviceId]
	if !ok {
		return nil, nil
	}

	return &d, nil
}

func (c *StorageMock) StoreSession(session mdl.Session) error {
	c.sessions[session.ID] = session
	return nil
}

func (c *StorageMock) FindSession(id string) (*mdl.Session, error) {
	session, ok := c.sessions[id]
	if !ok {
		return nil, errSessionNotFound
	}
	return &session, nil
}

func (c *StorageMock) HandlerLogin(session mdl.Session, ctx LoginContext) error {
	return nil
}
