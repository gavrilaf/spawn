package storage

import (
	"github.com/gavrilaf/go-auth/auth/cerr"
)

type MemorySessionStorage struct {
	Sessions map[string]Session
}

func NewMemorySessionsStorage() MemorySessionStorage {
	return MemorySessionStorage{Sessions: make(map[string]Session)}
}

func (p MemorySessionStorage) StoreSession(session Session) error {
	p.Sessions[session.ID] = session
	return nil
}

func (p MemorySessionStorage) FindSessionByID(id string) (*Session, error) {
	session, ok := p.Sessions[id]
	if !ok {
		return nil, cerr.SessionNotFound
	}
	return &session, nil
}
