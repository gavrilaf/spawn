package auth

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
