package go_auth

type User struct {
	ID        string
	Email     string
	Signature string
	Devices   map[string]string
}

type UserStorage interface {
	OpenStorage() error

	FindUser(email string) (*User, error)
}

/////////////////////////////////////////////////////////////

type ClientKey struct {
	ClientID string
	Secret   string
}

type ClientsStorage interface {
	FindClient(clientId string) (*ClientKey, error)
}

/////////////////////////////////////////////////////////////

type TokenStorage interface {
	IsTokenValid(accessToken string) bool
}

// Login form structure.
type Login struct {
	ClientID   string `form:"client_id" json:"client_id" binding:"required"`
	Username   string `form:"username" json:"username" binding:"required"`
	DeviceID   string `form:"device_id" json:"device_id" binding:"required"`
	SignSecret string `form:"sign_secret" json:"sign_secret" binding:"required"`
	SignKey    string `form:"sign_key" json:"sign_key" binding:"required"`
}
