package backend

import (
	"github.com/gavrilaf/spawn/pkg/env"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type BackendBridge struct {
	Client SpawnClient
	conn   *grpc.ClientConn
}

func (b *BackendBridge) Close() {
	b.conn.Close()
}

func CreateClient(en *env.Environment) (*BackendBridge, error) {
	log.Info("1")
	conn, err := grpc.Dial("localhost:7887", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	log.Info("2")

	client := NewSpawnClient(conn)
	log.Info("3")
	return &BackendBridge{client, conn}, nil
}
