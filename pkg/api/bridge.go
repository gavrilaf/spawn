package api

import (
	"context"
	"os"
	"time"

	"github.com/gavrilaf/spawn/pkg/cache"
	rdm "github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/gavrilaf/spawn/pkg/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Bridge struct {
	ReadModel  cache.Cache
	WriteModel *pb.BackendBridge
}

func CreateBridge(en *senv.Environment) *Bridge {
	log.Info("Starting Spawn api ...")

	log.Info("System environment:")
	for _, e := range os.Environ() {
		log.Info(e)
	}

	log.Infof("API environment type: %v", en.GetName())
	log.Infof("Backend path: %v", en.GetRPCOpts().URL)
	log.Infof("Cache path: %v", en.GetRedisOpts().URL)

	log.Info("Connecting to cache...")

	cache := cache.Connect(en)

	_, err := utils.Repeat(func() (interface{}, error) {
		err := cache.HealthCheck()
		if err != nil {
			log.Errorf("Cache healthcheck error: %v", err)
		}
		return nil, err
	}, 3, 3*time.Second)

	if err != nil {
		log.Errorf("Could not connect to cache: %v", err)
		return nil
	}

	log.Infof("Cache connection, ok")

	log.Info("Connecting to backend...")

	backend, err := utils.Repeat(func() (interface{}, error) {
		p, err := pb.CreateClient(en)
		if p != nil {
			_, err = p.Client.Ping(context.Background(), &pb.Empty{})
		}
		if err != nil {
			log.Errorf("Connect to the write model error: %v", err)
		}

		return p, err
	}, 3, 3*time.Second)

	if err != nil {
		log.Errorf("Could not connect to the write model: %v", err)
		return nil
	}

	log.Infof("Write model connection ok")

	return &Bridge{ReadModel: cache, WriteModel: backend.(*pb.BackendBridge)}
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
