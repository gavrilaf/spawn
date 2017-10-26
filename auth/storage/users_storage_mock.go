package storage

import (
	"fmt"
	"github.com/gavrilaf/go-auth/auth/cerr"
)

var Counter = 3

var DefaultUsers = []User{
	User{ID: "id1", Username: "id1@i.com", PasswordHash: "111111", Devices: []string{"d1", "d2"}},
	User{ID: "id2", Username: "id2@i.com", PasswordHash: "111111", Devices: []string{"d3", "d4"}},
}

type UsersStorageMock struct{}

func (p UsersStorageMock) AddUser(clientId string, deviceId string, username string, password string) error {
	if user, _ := p.FindUserByUsername(username); user != nil {
		return cerr.UserAlreadyExist
	}

	id := fmt.Sprintf("id%d", Counter)

	user := User{ID: id, Username: username, PasswordHash: password, Devices: []string{deviceId}}
	DefaultUsers = append(DefaultUsers, user)

	return nil
}

func (p UsersStorageMock) FindUserByUsername(username string) (*User, error) {
	for _, u := range DefaultUsers {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, cerr.UserUnknown
}
