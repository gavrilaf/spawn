package storage

type User struct {
	ID        string
	Email     string
	Signature string
	Devices   []string
}

type Session struct {
	SessionID string
	ClientID  string
	UserID    string
	Email     string
	Secret    string
}

type ClientKey struct {
	ClientID string
	Secret   string
}

/////////////////////////////////////////////////////////////////////////////////////

type ClientsStorage interface {
	FindClient(clientId string) (*ClientKey, error)
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
