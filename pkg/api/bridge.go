package api

import (
	"github.com/gavrilaf/spawn/pkg/cache"
	"github.com/gavrilaf/spawn/pkg/env"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
)

type Bridge struct {
	ReadModel  cache.Cache
	WriteModel *pb.BackendBridge
}

func CreateBridge(en *env.Environment) *Bridge {
	log.Infof("Starting storage with environment: %v", en.GetName())

	cache, err := cache.Connect(en)
	if err != nil {
		log.Errorf("Could not connect to the read  model: %v", err)
		return nil
	}
	log.Infof("Read model connection, ok")

	backend, err := pb.CreateClient(en)
	if err != nil {
		log.Errorf("Could not connect to the write model: %v", err)
		return nil
	}
	log.Infof("Write model connection, ok")

	return &Bridge{ReadModel: cache, WriteModel: backend}
}

func (p *Bridge) Close() {
	if p.ReadModel != nil {
		p.ReadModel.Close()
	}

	if p.WriteModel != nil {
		p.WriteModel.Close()
	}
}
