package types

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
	ClientID   string `json:"client_id" binding:"required"`
	Username   string `json:"username" binding:"required"`
	DeviceID   string `json:"device_id" binding:"required"`
	SignSecret string `json:"sign_secret" binding:"required"`
	SignKey    string `json:"sign_key" binding:"required"`
}
