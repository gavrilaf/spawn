package storage

type Client struct {
	ID     string
	Secret []byte
}

type User struct {
	ID        string
	Email     string
	Signature string
	Devices   []string
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
	FindUserByEmail(email string) (*User, error)
}

/////////////////////////////////////////////////////////////////////////////////////

type SessionsStorage interface {
	StoreSession(session Session) error
	FindSessionByID(id string) (*Session, error)
}
