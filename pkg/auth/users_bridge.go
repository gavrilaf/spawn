package auth

import (
	"fmt"
	"github.com/gavrilaf/spawn/pkg/env"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type UsersBridge struct {
	rpc *pb.BackendBridge
}

func CreateUsersBridge(en *env.Environment) *UsersBridge {
	rpc, err := pb.CreateClient(en)
	if err != nil {
		log.Errorf("Can not connect to backend: %v", err)
		panic(err)
	}
	log.Infof("Connected to the backend")

	return &UsersBridge{rpc}
}

func (b *UsersBridge) AddUser(clientId string, deviceId string, username string, password string) error {
	req := pb.RegisterUserRequest{
		Username:     username,
		PasswordHash: password,
		DeviceId:     deviceId}

	resp, err := b.rpc.Client.RegisterUser(context.Background(), &req)
	if resp != nil {
		log.Errorf("Registered user with id = %v", resp.Id)
	}

	return err
}

func (b *UsersBridge) FindUserByUsername(username string) (*User, error) {
	return nil, fmt.Errorf("not implemented")
}
