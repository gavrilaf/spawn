package senv

import (
	"os"
	"time"
)

func GetEnvironment() *Environment {
	name := getEnvVar("ENV_NAME", "Development")

	dbPath := getEnvVar("DATABASE_URL", "postgresql://spawnuser:spawn-pg-test-password@localhost/spawn?sslmode=disable")
	cachePath := getEnvVar("CACHE_URL", "redis://localhost:6379")
	queuePath := getEnvVar("QUEUE_URL", "amqp://localhost:5672")

	return &Environment{
		name: name,

		back: BackendOptions{
			URL:       queuePath,
			QueueName: "backend_queue",
			Timeout:   time.Duration(3) * time.Second,
		},

		cache: CacheOptions{
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
