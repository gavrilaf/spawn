package main

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/env"
)

type Cache struct {
}

func BuildCache() (*Cache, error) {
	e := env.GetEnvironment("")

	client := redis.NewClient(e.Redis())

	defer client.Close()

	pong, err := client.Ping().Result()
	log.Infof("%v, %v", pong, err)

	return &Cache{}, nil
}

func (p *Cache) PrintState() {

}
