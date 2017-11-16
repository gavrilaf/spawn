package main

import (
	"github.com/gavrilaf/spawn/pkg/env"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func newBackend() *BackendServer {
	en := env.GetEnvironment("Test")
	srv := CreateBackendServer(en)

	if srv == nil {
		panic("Can not create server")
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
