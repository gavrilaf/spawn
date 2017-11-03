package env

import (
	"github.com/go-redis/redis"
)

type Environment struct{}

func GetEnvironment(path string) *Environment {
	return &Environment{}
}

func (e *Environment) Redis() *redis.Options {
	return &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
}
