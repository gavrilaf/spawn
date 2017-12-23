package api

import (
	"github.com/gavrilaf/spawn/pkg/cache"
	rdm "github.com/gavrilaf/spawn/pkg/cache/model"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/gavrilaf/spawn/pkg/errx"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Bridge struct {
	ReadModel  cache.Cache
	WriteModel *pb.BackendBridge
}

func CreateBridge(en *env.Environment) *Bridge {
	log.Infof("Starting api with environment: %v", en.GetName())

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

///////////////////////////////////////////////////////////////////////////////
// Helpers

func (p *Bridge) HandleError(c *gin.Context, scope string, httpCode int, err error) {
	log.Errorf("api.Bridge.HandleError: scope=%v, code=%d, err=%v", scope, httpCode, err)
	errJSON := errx.Error2Map(err, scope)
	c.JSON(httpCode, gin.H{"error": errJSON})
	c.Abort()
}

func (p *Bridge) GetSession(c *gin.Context) (*rdm.Session, error) {
	return p.ReadModel.GetSession(c.GetString("session_id"))
}
