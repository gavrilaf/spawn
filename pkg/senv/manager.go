package senv

import (
	"os"
	"time"
)

func GetEnvironment() *Environment {
	name := getEnvVar("ENV_NAME", "Development")

	dbPath := getEnvVar("DB_PATH", "postgresql://spawnuser:spawn-pg-test-password@localhost/spawn?sslmode=disable")
	cachePath := getEnvVar("CACHE_PATH", "redis://localhost:6379")
	backendPath := getEnvVar("BACKEND_PATH", "localhost:7887")

	return &Environment{
		name: name,

		rpc: RPCOptions{
			URL:     backendPath,
			Timeout: time.Duration(3) * time.Second,
		},

		redis: RedisOptions{
			URL:         cachePath,
			MaxIdle:     3,
			IdleTimeout: time.Duration(240) * time.Second,
		},

		db: DBOptions{
			Driver:     "postgres",
			DataSource: dbPath,
		},
	}
}

func getEnvVar(key string, def string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return def
	}
	return v
}
