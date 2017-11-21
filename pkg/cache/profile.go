package cache

import (
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/model"
)

func profileRedisId(id string) string {
	return "profile:" + id
}

func (cache *Bridge) SetUserProfile(profile db.UserProfile) error {
	cacheProfile := mdl.CreateProfileFromDbModel(profile)
	_, err := cache.conn.Do("HMSET", redis.Args{}.Add(profileRedisId(cacheProfile.ID)).AddFlat(&cacheProfile)...)
	return err
}

func (cache *Bridge) GetUserProfile(userID string) (*mdl.UserProfile, error) {
	v, err := redis.Values(cache.conn.Do("HGETALL", profileRedisId(userID)))

	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	var profile mdl.UserProfile
	if err := redis.ScanStruct(v, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
