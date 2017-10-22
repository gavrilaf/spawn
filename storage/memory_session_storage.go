package storage

import (
	"fmt"
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
		return nil, fmt.Errorf("Cant find session by id %v", id)
	}
	return &session, nil
}
