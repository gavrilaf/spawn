package dbx

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"fmt"
	"github.com/gavrilaf/spawn/pkg/env"
	mdl "github.com/gavrilaf/spawn/pkg/model"
)

type Bridge struct {
	Db *sqlx.DB
}

func Connect(en *env.Environment) (*Bridge, error) {
	db, err := sqlx.Connect("postgres", "dbname=spawn host=localhost port=5432 user=postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to postgre database: %v", err)
	}

	return &Bridge{db}, nil
}

func (db *Bridge) Close() {
	if db.Db == nil {
		return
	}

	db.Db.Close()
}

///////////////////////////////////

type Database interface {
	RegisterUser(username string, password string, device mdl.DeviceInfo) (*mdl.UserProfile, error)

	UpdatePermissions(id string, permissons *mdl.Permissions) error
	UpdatePersonalInfo(id string, info *mdl.PersonalInfo) error

	GetUserProfile(id string) (*mdl.UserProfile, error)
	FindUserProfile(username string) (*mdl.UserProfile, error)

	AddDevice(userId string, device mdl.DeviceInfo) error
	ConfirmDevice(userId string, deviceId string) error
	RemoveDevice(userId string, deviceId string) error

	GetDevices(userId string) ([]mdl.DeviceInfo, error)
}
