package storage

import "encoding/hex"

type Client struct {
	id          string
	secret      string
	secretCache []byte
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Secret() []byte {
	if c.secretCache != nil {
		return c.secretCache
	}

	c.secretCache, _ = hex.DecodeString(c.secret)
	return c.secretCache
}

/////////////////////////////////////////////////////////////////////////////////////

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Devices      []string
}

type Session struct {
	ID           string
	RefreshToken string
	ClientID     string
	ClientSecret []byte
	UserID       string
}

/////////////////////////////////////////////////////////////////////////////////////

type ClientsStorage interface {
	FindClientByID(id string) (*Client, error)
}

/////////////////////////////////////////////////////////////////////////////////////

type UsersStorage interface {
	AddUser(clientId string, deviceId string, username string, password string) error
	FindUserByUsername(username string) (*User, error)
}

/////////////////////////////////////////////////////////////////////////////////////

type SessionsStorage interface {
	StoreSession(session Session) error
	FindSessionByID(id string) (*Session, error)
}
