package api

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/amqp/rpc"

	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/cache"
	rdm "github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/gavrilaf/spawn/pkg/utils"
)

type Bridge struct {
	ReadModel  cache.Cache
	WriteModel pb.SpawnClient
}

func CreateBridge(env *senv.Environment) *Bridge {
	log.Info("Starting Spawn api ...")

	log.Info("Connecting to cache...")

	cache := cache.Connect(env)

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
		srv, err := rpc.Connect(rpc.ClientConfig{
			Url:         env.GetBackOpts().URL,
			ServerQueue: env.GetBackOpts().QueueName,
			Timeout:     env.GetBackOpts().Timeout})

		if err != nil {
			log.Errorf("Could not connect to rpc: %v", err)
			return nil, err
		}

		client := pb.NewSpawnClient(srv)

		_, err = client.Ping(&pb.Empty{})
		if err != nil {
			client.Close()
			log.Errorf("Ping error: %v", err)
			return nil, err
		}

		return client, nil
	}, 3, 3*time.Second)

	if err != nil {
		log.Errorf("Could not connect to the write model: %v", err)
		return nil
	}

	log.Infof("Write model connection ok")

	return &Bridge{ReadModel: cache, WriteModel: backend.(pb.SpawnClient)}
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
