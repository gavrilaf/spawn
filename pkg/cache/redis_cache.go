package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gavrilaf/spawn/pkg/env"
)

type RedisCache struct {
	conn redis.Conn
}

func Connect(en *env.Environment) (*RedisCache, error) {

	//redis://[:password@]host[:port][/db-number][?option=value]
	conn, err := redis.DialURL("redis://localhost:7001")
	if err != nil {
		return nil, err
	}
	return &RedisCache{conn}, nil
}

func (cache *RedisCache) Close() error {
	return cache.conn.Close()
}
