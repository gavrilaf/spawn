package storage

import (
	"github.com/gavrilaf/go-auth/auth/cerr"
)

type MemorySessionStorage struct {
	sessions map[string]Session
}

func NewMemorySessionsStorage() *MemorySessionStorage {
	return &MemorySessionStorage{sessions: make(map[string]Session)}
}

func (p *MemorySessionStorage) StoreSession(session Session) error {
	p.sessions[session.ID] = session
	return nil
}

func (p *MemorySessionStorage) FindSessionByID(id string) (*Session, error) {
	session, ok := p.sessions[id]
	if !ok {
		return nil, cerr.SessionNotFound
	}
	return &session, nil
}
