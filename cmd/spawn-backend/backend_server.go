package main

import (
	"github.com/gavrilaf/spawn/pkg/cache"
	"github.com/gavrilaf/spawn/pkg/dbx"
	"github.com/gavrilaf/spawn/pkg/env"
	log "github.com/sirupsen/logrus"
)

type BackendServer struct {
	Db    *dbx.Bridge
	Cache *cache.Cache
}

func CreateBackendServer(en *env.Environment) *BackendServer {
	log.Infof("Starting backend with environment: %v", en.GetName())

	db, err := dbx.Connect(en)
	if err != nil {
		log.Errorf("Can not connect to database: %v", err)
		return nil
	}
	log.Infof("Db connection, ok")

	cache, err := cache.Connect(en)
	if err != nil {
		log.Errorf("Can not connect to cache: %v", err)
		return nil
	}
	log.Infof("Cache connection, ok")

	return &BackendServer{Db: db, Cache: cache}
}
