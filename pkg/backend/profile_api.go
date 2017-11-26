package backend

import (
	//"fmt"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (srv *Server) DoConfirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.Empty, error) {
	log.Infof("DoConfirm, %v", spew.Sdump(req))

	switch req.Kind {

	}

	return &pb.Empty{}, nil
}

func (srv *Server) UpdateUserPersonalInfo(ctx context.Context, req *pb.UserPersonalInfoRequest) (*pb.Empty, error) {
	return nil, fmt.Errorf("not implemented")
}

func (srv *Server) UpdateUserCountry(ctx context.Context, eq *pb.UserPersonalInfoRequest) (*pb.Empty, error) {
	return nil, fmt.Errorf("not implemented")
}
