package backend

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/gavrilaf/spawn/pkg/backend/pb"
)

func TestServer_AddDevice(t *testing.T) {
	srv := createTestSrv(t)
	defer srv.Close()

	userID := regTestUser(t, srv, testDevice)

	device := pb.Device{ID: "d2",
		Name:   "Test device",
		Locale: "ru",
		Lang:   "es"}

	_, err := srv.AddDevice(&pb.UserDevice{UserID: userID, Device: &device})
	assert.Nil(t, err)

	// New device is not confirmed; confirm code is added
	code, err := srv.cache.GetDeviceConfirmCode(userID, device.ID)
	assert.NotEmpty(t, code)
	assert.Nil(t, err)

	// New device added to the db
	dbDevices, err := srv.db.GetUserDevices(userID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(dbDevices))

	// and cache
	cacheDevice, err := srv.cache.GetDevice(userID, "d2")
	assert.Nil(t, err)
	assert.NotNil(t, cacheDevice)

	assert.Equal(t, "d2", cacheDevice.DeviceID)
	assert.Equal(t, userID, cacheDevice.UserID)
	assert.Equal(t, "ru", cacheDevice.Locale)
	assert.Equal(t, "es", cacheDevice.Lang)
	assert.Equal(t, false, cacheDevice.IsConfirmed)
}

func TestServer_ConfirmDevice(t *testing.T) {
	srv := createTestSrv(t)
	defer srv.Close()

	userID := regTestUser(t, srv, testDevice)
	deviceID := "d2"

	device := pb.Device{ID: deviceID,
		Name:   "Test device",
		Locale: "ru",
		Lang:   "es"}

	_, err := srv.AddDevice(&pb.UserDevice{UserID: userID, Device: &device})
	assert.Nil(t, err)

	code, err := srv.cache.GetDeviceConfirmCode(userID, device.ID)
	assert.NotEmpty(t, code)
	assert.Nil(t, err)

	sessionID := makeFakeSession(t, srv, userID, deviceID, false)
	_, err = srv.ConfirmDevice(&pb.ConfirmDeviceReq{SessionId: sessionID, Code: code})
	assert.Nil(t, err)

	// Should be confirmed in database
	db_device, err := srv.db.GetUserDevice(userID, deviceID)
	require.Nil(t, err)
	assert.True(t, db_device.IsConfirmed)

	// Should be session updated
	session, err := srv.cache.GetSession(sessionID)
	require.Nil(t, err)
	assert.True(t, session.IsDeviceConfirmed) // Device is confirmed now

	// Should be removed confrim code
	code, err = srv.cache.GetDeviceConfirmCode(userID, device.ID)
	assert.NotNil(t, err)
}

func TestServer_DeleteDevice(t *testing.T) {

}
