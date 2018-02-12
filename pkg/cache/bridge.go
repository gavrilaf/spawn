package cache

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gavrilaf/spawn/pkg/senv"
)

type Bridge struct {
	pool *redis.Pool
}

func (p *Bridge) Close() error {
	if p == nil || p.pool == nil {
		return nil
	}

	return p.pool.Close()
}

func (p *Bridge) HealthCheck() error {
	_, err := p.get().Do("PING")
	return err
}

//////////////////////////////////////////////////////////////////////////////////////////

func newPool(en *senv.Environment) *redis.Pool {
	return &redis.Pool{

		MaxIdle:     en.GetRedisOpts().MaxIdle,
		IdleTimeout: en.GetRedisOpts().IdleTimeout,

		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(en.GetRedisOpts().URL)
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

func deleteKeys(conn redis.Conn, pattern string) (int, error) {
	keys, err := getKeys(conn, pattern)
	if err != nil {
		return 0, err
	}

	for _, key := range keys {
		conn.Do("DEL", key)
	}

	return len(keys), nil
}
