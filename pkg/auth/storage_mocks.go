package auth

import (
	"fmt"
)

/*
 * Clients storage
 */

var DefaultClients = map[string]string{
	"client_test": "client_test_key",
	"client_ios":  "client_ios_key",
}

type ClientsStorageMock struct{}

func NewClientsStorageMock() *ClientsStorageMock {
	return &ClientsStorageMock{}
}

func (c *ClientsStorageMock) FindClientByID(id string) (*Client, error) {
	secret, ok := DefaultClients[id]
	if !ok {
		return nil, errClientUnknown
	}
	return &Client{id: id, secret: secret}, nil
}

/*
 * Users storage
 */

var DefaultUsers = []User{
	User{ID: "id1", Username: "id1@i.com", PasswordHash: "111111", Devices: []string{"d1", "d2"}},
	User{ID: "id2", Username: "id2@i.com", PasswordHash: "111111", Devices: []string{"d3", "d4"}},
}

type UsersStorageMock struct {
	registered []User
	counter    int
}

func NewUsersStorageMock() *UsersStorageMock {
	users := make([]User, len(DefaultUsers))
	for i, u := range DefaultUsers {
		users[i] = u
	}

	return &UsersStorageMock{registered: users, counter: 3}
}

func (p *UsersStorageMock) AddUser(clientId string, deviceId string, username string, password string) error {
	if user, _ := p.FindUserByUsername(username); user != nil {
		return errUserAlreadyExist
	}

	id := fmt.Sprintf("id%d", p.counter)
	p.counter += 1

	user := User{ID: id, Username: username, PasswordHash: password, Devices: []string{deviceId}}
	p.registered = append(p.registered, user)

	return nil
}

func (p *UsersStorageMock) FindUserByUsername(username string) (*User, error) {
	for _, u := range p.registered {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, errUserUnknown
}
