package backend

import (
	"github.com/gavrilaf/spawn/pkg/backend/pb"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/utils"

	log "github.com/sirupsen/logrus"
	"time"
)

func (srv *Server) UpdateUserPersonalInfo(arg *pb.UserPersonalInfo) (*pb.Empty, error) {
	go func() {
		log.Infof("UpdateUserPersonalInfo: %s", arg.String())

		tm := utils.CreateDate(int(arg.BirthDate.Year), time.Month(arg.BirthDate.Month), int(arg.BirthDate.Day))
		err := srv.db.UpdateUserPersonalInfo(arg.UserID, db.PersonalInfo{
			FirstName: arg.FirstName,
			LastName:  arg.LastName,
			BirthDate: tm,
		})

		if err != nil {
			log.Errorf("Could not update personal info in db: %v", err)
			return
		}

		srv.updateCachedUserProfile(arg.UserID)
	}()

	return &pb.Empty{}, nil
}

func (srv *Server) UpdateUserCountry(arg *pb.UserCountry) (*pb.Empty, error) {
	go func() {
		log.Infof("UpdateUserCountry: %s", arg.String())

		err := srv.db.UpdateUserCountry(arg.UserID, arg.Country)
		if err != nil {
			log.Errorf("Could not update country in db: %v", err)
			return
		}

		srv.updateCachedUserProfile(arg.UserID)
	}()

	return &pb.Empty{}, nil
}
