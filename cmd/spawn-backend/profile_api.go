package main

import (
	"fmt"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"golang.org/x/net/context"
)

func (src *BackendServer) RegisterUser(context.Context, *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {

	// Add user to the DB

	// Update Redis cache

	// Send email

	return nil, fmt.Errorf("not implemented")
}
