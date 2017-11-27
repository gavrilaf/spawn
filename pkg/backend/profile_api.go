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

func (srv *Server) UpdateUserCountry(ctx context.Context, req *pb.UserCountryRequest) (*pb.Empty, error) {
	go func() {
		log.Infof("UpdateUserCountry: %v", spew.Sdump(req))
		err := srv.db.UpdateUserCountry(req.UserID, req.Country)
		if err != nil {
			log.Errorf("Could not update country in db: %v", err)
			return
		}

		profile, err := srv.db.GetUserProfile(req.UserID)
		if err != nil {
			log.Errorf("Could not read profile from db: %v", err)
			return
		}

		err = srv.cache.SetUserProfile(*profile)
		if err != nil {
			log.Errorf("Could not read profile in cache: %v", err)
		}
	}()

	return &pb.Empty{}, nil
}
