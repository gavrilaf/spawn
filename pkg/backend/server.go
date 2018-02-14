package backend

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/gavrilaf/spawn/pkg/cache"
	"github.com/gavrilaf/spawn/pkg/dbx"
	"github.com/gavrilaf/spawn/pkg/errx"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/gavrilaf/spawn/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type ServerState int32

const (
	StateCreated ServerState = iota
	StateLoading
	StateOk
	StateError
)

type Server struct {
	db    dbx.Database
	cache cache.Cache
	state ServerState
	wg    *sync.WaitGroup
}

func CreateServer(en *senv.Environment) *Server {
	log.Info("Starting Spawn backened server...")
	log.Info("System environment:")
	for _, e := range os.Environ() {
		log.Info(e)
	}

	log.Infof("Backend environment type: %v", en.GetName())
	log.Infof("DB path: %v", en.GetDBOpts().DataSource)
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

	log.Info("Connecting to db...")

	db, err := utils.Repeat(func() (interface{}, error) {
		db, err := dbx.Connect(en)
		if err != nil {
			log.Errorf("Connect to db error: %v", err)
		}
		return db, err
	}, 3, 3*time.Second)

	if err != nil {
		log.Errorf("Could not connect to database: %v", err)
		return nil
	}

	log.Infof("Db connection, ok")

	return &Server{db: db.(dbx.Database), cache: cache, state: StateCreated, wg: &sync.WaitGroup{}}
}

func (srv *Server) StartServer() {
	log.Infof("Server started...")
	srv.updateServerState(StateLoading)

	srv.wg.Add(1)

	go srv.loadAuthCache()

	timeout := utils.WaitWithTimeout(srv.wg, 10*time.Second)

	if srv.state == StateLoading {
		if timeout {
			log.Errorf("Server loading timeout")
			srv.updateServerState(StateError)
		} else {
			srv.updateServerState(StateOk)
		}
	}

	log.Infof("Server initializing finished with state %d", srv.state)
}

func (srv *Server) GetServerState() ServerState {
	return srv.state
}

func (srv *Server) Close() {
	if srv.db != nil {
		err := srv.db.Close()
		log.Info("Closed database with err: %v", err)
	}

	if srv.cache != nil {
		err := srv.cache.Close()
		log.Info("Closed read model with err: %v", err)
	}

	srv.state = StateCreated
}

// API
func (srv *Server) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	state := srv.GetServerState()
	if state != StateOk {
		return nil, errx.ErrEnvironment(ErrScope, "Backed is not ready yet, current state: %d", state)
	}

	return &pb.Empty{}, nil
}

// Private
func (srv *Server) updateServerState(newState ServerState) {
	srv.state = newState
}
