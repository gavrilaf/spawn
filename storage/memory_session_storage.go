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
	p.Sessions[session.Token] = session
	return nil
}

func (p MemorySessionStorage) FindSessionByToken(token string) (*Session, error) {
	session, ok := p.Sessions[token]
	if !ok {
		return nil, fmt.Errorf("Cant find session by token %v", token)
	}
	return &session, nil
}
