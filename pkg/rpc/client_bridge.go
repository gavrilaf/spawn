package backend

import (
	"context"

	"github.com/gavrilaf/spawn/pkg/env"
	"google.golang.org/grpc"
)

type BackendBridge struct {
	Client     SpawnClient
	Ctx        context.Context
	CancelFunc context.CancelFunc
	conn       *grpc.ClientConn
}

func (b *BackendBridge) Close() {
	b.conn.Close()
}

func CreateClient(en *env.Environment) (*BackendBridge, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}

	ctx, cancel := context.WithTimeout(context.Background(), en.GetRPCOpts().Timeout)

	conn, err := grpc.DialContext(ctx, en.GetRPCOpts().URL, opts...)
	if err != nil {
		return nil, err
	}

	client := NewSpawnClient(conn)
	return &BackendBridge{client, ctx, cancel, conn}, nil
}
