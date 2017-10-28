package auth

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
