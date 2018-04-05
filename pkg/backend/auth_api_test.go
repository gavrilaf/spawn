package backend

import (
	"fmt"
	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/gavrilaf/spawn/pkg/backend/pb"
	"github.com/gavrilaf/spawn/pkg/senv"
)

func createServer(t *testing.T) *Server {
	en := senv.GetEnvironment()
	srv := CreateServer(en)
	require.NotNil(t, srv)
	return srv
}

func TestServer_RegisterUser(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()

	username := uuid.NewV4().String() + "@spawn.com"

	arg := pb.CreateUserReq{
		Username:     username,
		PasswordHash: "123456",
		Device: &pb.Device{
			ID:     "device-1",
			Name:   "Test device",
			Locale: "ru",
			Lang:   "es"},
	}

	res, err := srv.CreateUser(&arg)
	assert.Nil(t, err)

	dbUser, err := srv.db.GetUserProfile(res.ID)
	assert.Nil(t, err)
	assert.NotNil(t, dbUser)
	assert.Equal(t, username, dbUser.Username)

	dbDevices, err := srv.db.GetUserDevices(res.ID)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(dbDevices))
	assert.Equal(t, "device-1", dbDevices[0].ID)
	assert.Equal(t, "Test device", dbDevices[0].Name)
	assert.Equal(t, dbUser.ID, dbDevices[0].UserID)
	assert.Equal(t, "ru", dbDevices[0].Locale)
	assert.Equal(t, "es", dbDevices[0].Lang)
	assert.Equal(t, true, dbDevices[0].IsConfirmed)

	cacheUser, err := srv.cache.FindUserAuthInfo(username)
	assert.Nil(t, err)
	assert.NotNil(t, cacheUser)
	assert.Equal(t, username, cacheUser.Username)
	assert.Equal(t, dbUser.ID, cacheUser.ID)

	cacheDevice, err := srv.cache.GetDevice(res.ID, "device-1")
	assert.Nil(t, err)
	assert.NotNil(t, cacheDevice)

	assert.Equal(t, "device-1", cacheDevice.DeviceID)
	assert.Equal(t, dbUser.ID, cacheDevice.UserID)
	assert.Equal(t, "ru", cacheDevice.Locale)
	assert.Equal(t, "es", cacheDevice.Lang)
	assert.Equal(t, true, cacheDevice.IsConfirmed)
}

func TestServer_AddDevice(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()

	username := uuid.NewV4().String() + "@spawn.com"

	device := pb.Device{
		ID:     "d1",
		Name:   "Test device",
		Locale: "ru",
		Lang:   "ru"}

	userArg := pb.CreateUserReq{
		Username:     username,
		PasswordHash: "123456",
		Device:       &device,
	}

	user, err := srv.CreateUser(&userArg)
	assert.Nil(t, err)

	device.ID = "d2"
	_, err = srv.AddDevice(&pb.UserDevice{UserID: user.ID, Device: &device})
	assert.Nil(t, err)

	code, err := srv.cache.GetDeviceConfirmCode(user.ID, device.ID)
	assert.NotEmpty(t, code)
	assert.Nil(t, err)

	fmt.Printf("Confirmation code is %v\n", code)

	dbDevices, err := srv.db.GetUserDevices(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(dbDevices))

	cacheDevice, err := srv.cache.GetDevice(user.ID, "d2")
	assert.Nil(t, err)
	assert.NotNil(t, cacheDevice)

	assert.Equal(t, "d2", cacheDevice.DeviceID)
	assert.Equal(t, user.ID, cacheDevice.UserID)
	assert.Equal(t, "ru", cacheDevice.Locale)
	assert.Equal(t, "ru", cacheDevice.Lang)
	assert.Equal(t, false, cacheDevice.IsConfirmed)
}

func TestServer_HandleLogin(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()

	username := uuid.NewV4().String() + "@spawn.com"

	device := pb.Device{
		ID:     "d1",
		Name:   "Test device",
		Locale: "ru",
		Lang:   "ru"}

	userID, err := srv.CreateUser(&pb.CreateUserReq{
		Username:     username,
		PasswordHash: "123456",
		Device:       &device,
	})
	assert.Nil(t, err)

	devices, err := srv.db.GetUserDevicesEx(userID.ID)
	require.Nil(t, err)
	assert.Equal(t, 1, len(devices))

	assert.Nil(t, devices[0].GetLoginTime())
	assert.Equal(t, "", devices[0].GetUserAgent())
	assert.Equal(t, "", devices[0].GetLoginIP())
	assert.Equal(t, "", devices[0].GetLoginRegion())
	assert.Equal(t, "ru", devices[0].Locale)
	assert.Equal(t, "ru", devices[0].Lang)
	assert.Equal(t, "Test device", devices[0].Name)

	device.Lang = "es"
	device.Locale = "en"
	device.Name = "Updated"

	loginArg := pb.LoginReq{SessionID: "", UserID: userID.ID, Device: &device, UserAgent: "test", LoginIP: "127.0.0.1", LoginRegion: "ru"}

	_, err = srv.HandleLogin(&loginArg)
	assert.Nil(t, err)

	devices, err = srv.db.GetUserDevicesEx(userID.ID)
	require.Nil(t, err)
	assert.Equal(t, 1, len(devices))
	assert.Equal(t, "test", devices[0].GetUserAgent())
	assert.Equal(t, "127.0.0.1", devices[0].GetLoginIP())
	assert.Equal(t, "ru", devices[0].GetLoginRegion())
	assert.Equal(t, "en", devices[0].Locale)
	assert.Equal(t, "es", devices[0].Lang)
	assert.Equal(t, "Updated", devices[0].Name)
}
