package storage

type StorageFacade struct {
	Clients  ClientsStorage
	Users    UsersStorage
	Sessions SessionsStorage
}

func (p *StorageFacade) FindClient(clientId string) (*ClientKey, error) {
	return p.Clients.FindClient(clientId)
}

func (p *StorageFacade) FindUserByEmail(email string) (*User, error) {
	return p.Users.FindUserByEmail(email)
}

func (p *StorageFacade) StoreSession(session Session) error {
	return p.Sessions.StoreSession(session)
}

func (p *StorageFacade) FindSessionByID(id string) (*Session, error) {
	return p.Sessions.FindSessionByID(id)
}
