package cache

import (
	"github.com/garyburd/redigo/redis"
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func profileID(id string) string {
	return "profile:" + id
}

func (br *Bridge) SetUserProfile(profile db.UserProfile) error {
	conn := br.get()
	defer conn.Close()

	cacheProfile := mdl.CreateProfileFromDbModel(profile)
	_, err := conn.Do("HMSET", redis.Args{}.Add(profileID(cacheProfile.ID)).AddFlat(&cacheProfile)...)
	return err
}

func (br *Bridge) GetUserProfile(userID string) (*mdl.UserProfile, error) {
	conn := br.get()
	defer conn.Close()

	key := profileID(userID)
	v, err := redis.Values(conn.Do("HGETALL", key))

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
