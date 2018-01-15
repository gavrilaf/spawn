package senv

import (
	"time"
)

func getTestEnv() *Environment {
	return &Environment{
		name: "Test",

		rpc: RPCOptions{
			URL:     "localhost:7887",
			Timeout: time.Duration(3) * time.Second,
		},

		redis: RedisOptions{
			URL:         "redis://localhost:6379",
			MaxIdle:     3,
			IdleTimeout: time.Duration(240) * time.Second,
		},

		db: DBOptions{
			Driver:     "postgres",
			DataSource: "dbname=spawn host=localhost port=5432 user=spawnuser password=spawn-pg-test-password sslmode=disable",
		},
	}
}
