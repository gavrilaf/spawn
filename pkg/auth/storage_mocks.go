package auth

import (
	"fmt"
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

var DefaultUsers = []mdl.UserProfile{
	mdl.UserProfile{
		ID: "id1",
		AuthInfo: mdl.AuthInfo{
			Username:         "id1@spawn.com",
			PasswordHash:     "24326124313024636d636548593741794976623167777154732f502f754c4c466a4c755543784a6b696f386f6b4c344a42686a514b76494943654653",
			IsLocked:         false,
			IsEmailConfirmed: false,
			Is2FARequired:    false},
		PersonalInfo: mdl.PersonalInfo{
			FirstName: "FirstName",
			LastName:  "LastName"}},

	mdl.UserProfile{
		ID: "id2",
		AuthInfo: mdl.AuthInfo{
			Username:         "id2@spawn.com",
			PasswordHash:     "243261243130247a6b4576684654664b4c4a5945486871766e6b51472e354771335676664a492f6a5232304c73465774354c553939314b3944766b6d",
			IsLocked:         false,
			IsEmailConfirmed: false,
			Is2FARequired:    false},
		PersonalInfo: mdl.PersonalInfo{
			FirstName: "FirstName",
			LastName:  "LastName"}},
}

var DefaultDevices = map[string]map[string]bool{
	"id1": {"d1": true, "d2": true},
	"id2": {"d3": true, "d4": true},
}

type StorageMock struct {
	users    map[string]mdl.UserProfile
	counter  int
	devices  map[string]map[string]bool
	sessions map[string]mdl.Session
}

func NewStorageMock(en *env.Environment) *StorageMock {
	users := make(map[string]mdl.UserProfile)

	for _, v := range DefaultUsers {
		users[v.Username] = v
	}

	return &StorageMock{users: users, counter: 3, devices: DefaultDevices, sessions: make(map[string]mdl.Session)}
}

///////

func (c *StorageMock) FindClient(id string) (*mdl.Client, error) {
	secret, ok := DefaultClients[id]
	if !ok {
		return nil, errClientUnknown
	}
	return &mdl.Client{ID: id, Secret: secret}, nil
}

func (c *StorageMock) RegisterUser(username string, password string, deviceId string) error {
	if user, _ := c.FindUser(username); user != nil {
		return errUserAlreadyExist
	}

	id := fmt.Sprintf("id%d", c.counter)
	c.counter += 1
	user := mdl.UserProfile{
		ID: id,
		AuthInfo: mdl.AuthInfo{
			Username:         username,
			PasswordHash:     password,
			IsLocked:         false,
			IsEmailConfirmed: false,
			Is2FARequired:    false}}

	c.users[username] = user
	c.devices[id] = map[string]bool{deviceId: true}

	return nil
}

func (c *StorageMock) FindUser(username string) (*mdl.UserProfile, error) {
	profile, ok := c.users[username]
	if !ok {
		return nil, errUserUnknown
	}

	return &profile, nil
}

func (c *StorageMock) IsDeviceAllowed(userId string, deviceId string) (bool, error) {
	devices, ok := c.devices[userId]
	if ok {
		return devices[deviceId], nil
	}
	return false, nil
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
