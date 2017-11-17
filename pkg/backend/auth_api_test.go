package backend

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gavrilaf/spawn/pkg/env"
	//mdl "github.com/gavrilaf/spawn/pkg/model"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/gavrilaf/spawn/pkg/rpc"
	"golang.org/x/net/context"
)

func createBackend(t *testing.T) *Server {
	en := env.GetEnvironment("Test")
	srv := CreateServer(en)
	require.NotNil(t, srv)
	return srv
}

func TestRegisterUser(t *testing.T) {

	srv := createBackend(t)

	username := uuid.NewV4().String() + "@spawn.com"

	req := pb.CreateUserRequest{
		Username:     username,
		PasswordHash: "123456",
		Device: &pb.Device{
			Id:     "device-1",
			Name:   "Test device",
			Locale: "ru",
			Lang:   "ru"},
	}

	res, err := srv.CreateUser(context.Background(), &req)
	assert.Nil(t, err)

	fmt.Printf("Registered user: %v\n", spew.Sdump(res))

	dbUser, err := srv.Db.GetUserProfile(res.UserId)
	assert.Nil(t, err)
	assert.NotNil(t, dbUser)

	fmt.Printf("Db user: %v\n", spew.Sdump(dbUser))

	dbDevices, err := srv.Db.GetDevices(res.UserId)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(dbDevices))
	assert.Equal(t, true, dbDevices[0].IsConfirmed)

	fmt.Printf("Db devices: %v\n", spew.Sdump(dbDevices))

	cacheUser, err := srv.Cache.FindUserAuthInfo(username)
	assert.Nil(t, err)
	assert.NotNil(t, cacheUser)

	fmt.Printf("Cache user: %v\n", spew.Sdump(cacheUser))

	cacheDevice, err := srv.Cache.FindDevice(res.UserId, "device-1")
	assert.Nil(t, err)
	assert.NotNil(t, cacheDevice)

	assert.Equal(t, true, cacheDevice.IsConfirmed)

	fmt.Printf("Cache device: %v\n", spew.Sdump(cacheDevice))

}
