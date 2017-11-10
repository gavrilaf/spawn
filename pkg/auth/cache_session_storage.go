package auth

import (
	"github.com/gavrilaf/spawn/pkg/cache"
	"github.com/gavrilaf/spawn/pkg/env"
	log "github.com/sirupsen/logrus"
)

type CacheSessionStorage struct {
	rc *cache.RedisCache
}

func NewCacheSessionsStorage(e *env.Environment) *CacheSessionStorage {
	p, err := cache.Connect(e)
	if err != nil {
		log.Errorf("Can not connect to cache: %v", err)
		panic(err)
	}
	log.Infof("Sessions Redis storage connected")

	return &CacheSessionStorage{p}
}

func (p *CacheSessionStorage) StoreSession(session Session) error {

	sn := cache.Session{
		ID:           session.ID,
		RefreshToken: session.RefreshToken,
		ClientID:     session.ClientID,
		ClientSecret: session.ClientSecret,
		UserID:       session.UserID,
		DeviceID:     session.DeviceID}

	return p.rc.AddSession(sn)
}

func (p *CacheSessionStorage) FindSessionByID(id string) (*Session, error) {
	sn, err := p.rc.FindSession(id)
	if err != nil {
		log.Errorf("Can't find session with id %v, error: %v", id, err)
		return nil, errSessionNotFound
	}

	session := Session{
		ID:           sn.ID,
		RefreshToken: sn.RefreshToken,
		ClientID:     sn.ClientID,
		ClientSecret: sn.ClientSecret,
		UserID:       sn.UserID,
		DeviceID:     sn.DeviceID}

	return &session, nil
}
