package auth

type Client struct {
	id     string
	secret string
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Secret() []byte {
	return []byte(c.secret)
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
	DeviceID     string
}
