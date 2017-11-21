package dbx

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver module

	"fmt"

	mdl "github.com/gavrilaf/spawn/pkg/dbx/model"
	"github.com/gavrilaf/spawn/pkg/env"
)

type Bridge struct {
	conn *sqlx.DB
}

func Connect(en *env.Environment) (Database, error) {
	db, err := sqlx.Connect("postgres", "dbname=spawn host=localhost port=5432 user=postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to postgre database: %v", err)
	}

	return &Bridge{db}, nil
}

// Database is an interface to the Spawn storage
type Database interface {
	Close()

	RegisterUser(username string, password string, device mdl.DeviceInfo) (*mdl.UserProfile, error)

	GetUserProfile(id string) (*mdl.UserProfile, error)
	FindUserProfile(username string) (*mdl.UserProfile, error)

	UpdateUserPermissions(id string, permissions mdl.Permissions) error

	UpdateUserPersonalInfo(id string, info mdl.PersonalInfo) error
	UpdateUserCountry(id string, contry string) error
	UpdateUserPhoneNumber(id string, phone mdl.PhoneNumber) error

	ReadAllUserProfiles() (<-chan *mdl.UserProfile, <-chan error)

	AddDevice(userID string, device mdl.DeviceInfo) error
	RemoveDevice(userID string, deviceID string) error

	ConfirmDevice(userID string, deviceID string) error
	SetDeviceFingerprint(userID string, deviceID string, fingerprint []byte) error

	GetUserDevice(userID string, deviceID string) (*mdl.DeviceInfo, error)
	GetUserDevices(userID string) ([]mdl.DeviceInfo, error)
	GetUserDevicesEx(userID string) ([]mdl.DeviceInfoEx, error)

	LogUserLogin(userID string, deviceID string, ip string, region string) error
}

//////////////////////////////////////////////////////////////////////////////

func (db *Bridge) Close() {
	if db.conn == nil {
		return
	}

	db.conn.Close()
}
