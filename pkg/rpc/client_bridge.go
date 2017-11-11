package backend

import (
	"github.com/gavrilaf/spawn/pkg/env"
	"google.golang.org/grpc"
	"time"
)

type BackendBridge struct {
	Client SpawnClient
	conn   *grpc.ClientConn
}

func (b *BackendBridge) Close() {
	b.conn.Close()
}

func CreateClient(en *env.Environment) (*BackendBridge, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Duration(3) * time.Second),
	}

	conn, err := grpc.Dial("localhost:7887", opts...)
	if err != nil {
		return nil, err
	}

	client := NewSpawnClient(conn)
	return &BackendBridge{client, conn}, nil
}
