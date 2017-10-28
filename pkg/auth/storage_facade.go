package auth

type StorageFacade struct {
	Clients  ClientsStorage
	Users    UsersStorage
	Sessions SessionsStorage
}

func (p *StorageFacade) FindClientByID(id string) (*Client, error) {
	return p.Clients.FindClientByID(id)
}

func (p *StorageFacade) AddUser(clientId string, deviceId string, username string, password string) error {
	return p.Users.AddUser(clientId, deviceId, username, password)
}

func (p *StorageFacade) FindUserByUsername(username string) (*User, error) {
	return p.Users.FindUserByUsername(username)
}

func (p *StorageFacade) StoreSession(session Session) error {
	return p.Sessions.StoreSession(session)
}

func (p *StorageFacade) FindSessionByID(id string) (*Session, error) {
	return p.Sessions.FindSessionByID(id)
}
