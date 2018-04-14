package backend

import (
	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/gavrilaf/spawn/pkg/backend/pb"
)

func TestServer_RegisterUser(t *testing.T) {
	srv := createTestSrv(t)
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

	// User profile is saved to the db
	dbUser, err := srv.db.GetUserProfile(res.ID)
	assert.Nil(t, err)
	assert.NotNil(t, dbUser)
	assert.Equal(t, username, dbUser.Username)

	// Device is added
	dbDevices, err := srv.db.GetUserDevices(res.ID)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(dbDevices))
	assert.Equal(t, "device-1", dbDevices[0].DeviceID)
	assert.Equal(t, "Test device", dbDevices[0].Name)
	assert.Equal(t, dbUser.ID, dbDevices[0].UserID)
	assert.Equal(t, "ru", dbDevices[0].Locale)
	assert.Equal(t, "es", dbDevices[0].Lang)
	assert.Equal(t, true, dbDevices[0].IsConfirmed)

	// User auth info stores to the cache
	cacheUser, err := srv.cache.FindUserAuthInfo(username)
	assert.Nil(t, err)
	assert.NotNil(t, cacheUser)
	assert.Equal(t, username, cacheUser.Username)
	assert.Equal(t, dbUser.ID, cacheUser.ID)

	// Device info stores to the cache
	cacheDevice, err := srv.cache.GetDevice(res.ID, "device-1")
	assert.Nil(t, err)
	assert.NotNil(t, cacheDevice)

	assert.Equal(t, "device-1", cacheDevice.DeviceID)
	assert.Equal(t, dbUser.ID, cacheDevice.UserID)
	assert.Equal(t, "ru", cacheDevice.Locale)
	assert.Equal(t, "es", cacheDevice.Lang)
	assert.Equal(t, true, cacheDevice.IsConfirmed)
}

func TestServer_HandleLogin(t *testing.T) {
	srv := createTestSrv(t)
	defer srv.Close()

	userID := regTestUser(t, srv, testDevice)

	// After registration user sign in log is empty
	devices, err := srv.db.GetUserDevicesEx(userID)
	require.Nil(t, err)

	assert.Equal(t, int64(0), devices[0].GetLoginTime().Unix())
	assert.Equal(t, "", devices[0].GetUserAgent())
	assert.Equal(t, "", devices[0].GetLoginIP())
	assert.Equal(t, "", devices[0].GetLoginRegion())
	assert.Equal(t, "ru", devices[0].Locale)
	assert.Equal(t, "ru", devices[0].Lang)
	assert.Equal(t, "Test device", devices[0].Name)

	device := pb.Device{
		ID:     testDevice.ID,
		Name:   "Updated",
		Locale: "en",
		Lang:   "es"}

	loginArg := pb.LoginReq{SessionID: "", UserID: userID, Device: &device, UserAgent: "test", LoginIP: "127.0.0.1", LoginRegion: "ru"}

	// Handle user sign in
	_, err = srv.HandleLogin(&loginArg)
	assert.Nil(t, err)

	// User sign in log is updated
	devices, err = srv.db.GetUserDevicesEx(userID)
	require.Nil(t, err)
	assert.Equal(t, 1, len(devices))
	assert.Equal(t, "test", devices[0].GetUserAgent())
	assert.Equal(t, "127.0.0.1", devices[0].GetLoginIP())
	assert.Equal(t, "ru", devices[0].GetLoginRegion())
	assert.Equal(t, "en", devices[0].Locale)
	assert.Equal(t, "es", devices[0].Lang)
	assert.Equal(t, "Updated", devices[0].Name)
}
