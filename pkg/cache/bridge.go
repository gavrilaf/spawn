package cache

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gavrilaf/spawn/pkg/env"
)

type Bridge struct {
	pool *redis.Pool
}

func newPool(en *env.Environment) *redis.Pool {
	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://localhost:7001")
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (p *Bridge) Close() error {
	if p == nil || p.pool == nil {
		return nil
	}

	return p.pool.Close()
}

func (p *Bridge) get() redis.Conn {
	return p.pool.Get()
}

//////////////////////////////////////////////////////////////////////////////////////////

func getKeys(conn redis.Conn, pattern string) ([]string, error) {
	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}
