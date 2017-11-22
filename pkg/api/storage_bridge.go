package api

import (
	"github.com/gavrilaf/spawn/pkg/cache"
	"github.com/gavrilaf/spawn/pkg/env"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
)

type Bridge struct {
	Cache   cache.Cache
	Backend *pb.BackendBridge
}

func CreateBridge(en *env.Environment) *Bridge {
	log.Infof("Starting storage with environment: %v", en.GetName())

	cache, err := cache.Connect(en)
	if err != nil {
		log.Errorf("Can not connect to cache: %v", err)
		return nil
	}
	log.Infof("Cache connection, ok")

	backend, err := pb.CreateClient(en)
	if err != nil {
		log.Errorf("Can not connect to backend: %v", err)
		return nil
	}
	log.Infof("Backend connection, ok")

	return &Bridge{Cache: cache, Backend: backend}
}

func (p *Bridge) Close() {
	if p.Cache != nil {
		p.Cache.Close()
	}

	if p.Backend != nil {
		p.Backend.Close()
	}
}
