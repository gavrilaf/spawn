package cache

import (
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func profileRedisID(id string) string {
	return "profile:" + id
}

func (cache *Bridge) SetUserProfile(profile db.UserProfile) error {
	cacheProfile := mdl.CreateProfileFromDbModel(profile)
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(profileRedisID(cacheProfile.ID)).AddFlat(&cacheProfile)...)
	return err
}

func (cache *Bridge) GetUserProfile(userID string) (*mdl.UserProfile, error) {
	key := profileRedisID(userID)
	v, err := redis.Values(cache.conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, key)
	}

	var profile mdl.UserProfile
	if err := redis.ScanStruct(v, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
