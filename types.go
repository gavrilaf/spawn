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
