package main

import (
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func newServer() *BackendServer {
	s := new(BackendServer)
	return s
}

func main() {

	log.Info("Spawn backend started")

	lis, err := net.Listen("tcp", "localhost:7887")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSpawnServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}
