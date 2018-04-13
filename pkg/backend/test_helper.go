package backend

import (
	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/require"
	"testing"

	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/senv"
)

var testDevice = &pb.Device{
	ID:     "device-1",
	Name:   "Test device",
	Locale: "ru",
	Lang:   "ru"}

func createTestSrv(t *testing.T) *Server {
	en := senv.GetEnvironment()
	srv := CreateServer(en)
	require.NotNil(t, srv)
	return srv
}

func regTestUser(t *testing.T, srv *Server, device *pb.Device) string {
	username := uuid.NewV4().String() + "@spawn.com"
	arg := pb.CreateUserReq{
		Username:     username,
		PasswordHash: "123456",
		Device:       device,
	}

	res, err := srv.CreateUser(&arg)
	require.Nil(t, err)

	return res.ID
}

func makeFakeSession(t *testing.T, srv *Server, userID string, deviceID string, confirmed bool) string {
	session := mdl.Session{
		RefreshToken:      "",
		Nonce:             1,
		ClientID:          "",
		ClientSecret:      []byte(""),
		UserID:            userID,
		DeviceID:          deviceID,
		IsDeviceConfirmed: confirmed,
		Locale:            "",
		Lang:              "",
		Permissions:       db.Permissions{false, false, false, 0},
	}

	sessionID, err := srv.cache.AddSession(session, true)
	require.Nil(t, err)

	return sessionID
}
