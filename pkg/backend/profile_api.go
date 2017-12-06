package backend

import (
	"github.com/davecgh/go-spew/spew"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"github.com/golang/protobuf/ptypes"
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
	go func() {
		log.Infof("UpdateUserPersonalInfo: %v", spew.Sdump(req))

		tm, err := ptypes.Timestamp(req.BirthDate)
		if err != nil {
			log.Errorf("Invalid birth date: %v", err)
			return
		}

		err = srv.db.UpdateUserPersonalInfo(req.UserID, db.PersonalInfo{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			BirthDate: tm,
		})

		if err != nil {
			log.Errorf("Could not update personal info in db: %v", err)
			return
		}

		srv.updateCachedUserProfile(req.UserID)
	}()

	return &pb.Empty{}, nil

}

func (srv *Server) UpdateUserCountry(ctx context.Context, req *pb.UserCountryRequest) (*pb.Empty, error) {
	go func() {
		log.Infof("UpdateUserCountry: %v", spew.Sdump(req))
		err := srv.db.UpdateUserCountry(req.UserID, req.Country)
		if err != nil {
			log.Errorf("Could not update country in db: %v", err)
			return
		}

		srv.updateCachedUserProfile(req.UserID)
	}()

	return &pb.Empty{}, nil
}
