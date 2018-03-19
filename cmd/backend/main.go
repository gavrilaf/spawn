package main

import (
	"os"
	"time"

	"github.com/gavrilaf/amqp/rpc"
	log "github.com/sirupsen/logrus"

	"github.com/gavrilaf/spawn/pkg/backend"
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/gavrilaf/spawn/pkg/utils"
)

func newBackend(env *senv.Environment) *backend.Server {
	handler := backend.CreateServer(env)

	if handler == nil {
		log.Fatal("Could not create server")
	}

	handler.StartServer()

	if handler.GetServerState() != backend.StateOk {
		log.Fatal("Could not start server")
	}

	return handler
}

func main() {

	log.Info("Spawn backend starting...")

	env := senv.GetEnvironment()

	log.Info("System environment:")
	for _, e := range os.Environ() {
		log.Info(e)
	}

	log.Infof("Backend environment: %s", env.String())

	p, err := utils.Repeat(func() (interface{}, error) {
		srv, err := rpc.CreateServer(env.GetBackOpts().URL, env.GetBackOpts().QueueName)
		if err != nil {
			log.Errorf("%v", err)
		}
		return srv, err
	}, 3, 3*time.Second)

	if err != nil {
		log.Fatalf("Failed to create RPC server: %v", err)
	}

	srv := p.(rpc.Server)
	defer srv.Close()

	backend := newBackend(env)
	defer backend.Close()

	log.Infof("Run backend queue listener")
	pb.RunServer(srv, backend)
}
