package storage

type User struct {
	ID        string
	Email     string
	Signature string
	Devices   []string
}

type Session struct {
	Token     string
	SessionID string
	UserID    string
	Email     string
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
	FindSessionByToken(token string) (*Session, error)
}
