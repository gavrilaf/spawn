package storage

type Client struct {
	ID     string
	Secret string
}

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
	ClientSecret string
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
