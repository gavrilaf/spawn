package main

import (
	"net"

	"github.com/gavrilaf/spawn/pkg/backend"
	"github.com/gavrilaf/spawn/pkg/env"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func newBackend() *backend.Server {
	en := env.GetEnvironment("Test")
	if en == nil {
		panic("Could not read environment")
	}

	srv := backend.CreateServer(en)

	if srv == nil {
		panic("Could not create server")
	}

	srv.StartServer()

	if srv.GetServerState() != backend.StateOk {
		panic("Could not start server")
	}

	return srv
}

func main() {

	log.Info("Spawn backend starting...")

	lis, err := net.Listen("tcp", "localhost:7887")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	backServer := newBackend()

	pb.RegisterSpawnServer(grpcServer, backServer)

	grpcServer.Serve(lis)
}