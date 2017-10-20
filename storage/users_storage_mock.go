package storage

import (
	"fmt"
)

var DefaultUsers = []User{
	User{ID: "id1", Email: "id1@i.com", Signature: "111111", Devices: []string{"d1, d2"}},
	User{ID: "id2", Email: "id2@i.com", Signature: "111111", Devices: []string{"d3, d4"}},
}

type UsersStorageMock struct{}

func (p UsersStorageMock) FindUserByEmail(email string) (*User, error) {
	for _, u := range DefaultUsers {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("Email %v not found", email)
}
