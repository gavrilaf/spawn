package dbx

import (
	//"fmt"
	"github.com/gavrilaf/spawn/pkg/env"
	mdl "github.com/gavrilaf/spawn/pkg/model"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func GetBridge(t *testing.T) *Bridge {
	db, _ := Connect(env.GetEnvironment("Test"))
	require.NotNil(t, db)
	require.NotNil(t, db.Db)

	return db
}

func TestRegisterProfile(t *testing.T) {
	db := GetBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	device := mdl.DeviceInfo{ID: "d1", Name: "d1-name", IsConfirmed: true, LoginTime: time.Now()}

	profile, err := db.RegisterUser(username, "password", device)
	require.Nil(t, err)
	require.NotNil(t, profile)

	assert.Equal(t, false, profile.IsEmailConfirmed)
	assert.Equal(t, false, profile.Is2FARequired)
	assert.Equal(t, false, profile.IsLocked)

	p1, err := db.GetUserProfile(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, p1)

	p2, err := db.FindUserProfile(username)
	require.Nil(t, err)
	require.NotNil(t, p2)

	assert.Equal(t, profile, p1)
	assert.Equal(t, profile, p2)
	assert.Equal(t, p1, p2)

	devices, err := db.GetDevices(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, devices)

	assert.Equal(t, 1, len(devices))
	assert.Equal(t, devices[0].ID, device.ID)
	assert.Equal(t, devices[0].UserID, profile.ID)
}

func TestDeviceManagement(t *testing.T) {
	db := GetBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	d1 := mdl.DeviceInfo{ID: "d1", Name: "d1-name", IsConfirmed: true, LoginTime: time.Now()}
	d2 := mdl.DeviceInfo{ID: "d2", Name: "d2-name", IsConfirmed: false, LoginTime: time.Now()}
	d3 := mdl.DeviceInfo{ID: "d3", Name: "d3-name", IsConfirmed: false, LoginTime: time.Now()}

	profile, err := db.RegisterUser(username, "password", d1)
	require.Nil(t, err)

	devices, err := db.GetDevices(profile.ID)
	require.Nil(t, err)

	assert.Equal(t, 1, len(devices))
	assert.Equal(t, devices[0].IsConfirmed, true)

	err = db.AddDevice(profile.ID, d2)
	require.Nil(t, err)

	err = db.AddDevice(profile.ID, d3)
	require.Nil(t, err)

	devices, _ = db.GetDevices(profile.ID)

	assert.Equal(t, 3, len(devices))
	for _, d := range devices {
		if d.ID == "d2" || d.ID == "d3" {
			assert.Equal(t, false, d.IsConfirmed)
		}
	}

	err = db.ConfirmDevice(profile.ID, "d2")
	assert.Nil(t, err)

	err = db.ConfirmDevice(profile.ID, "d3")
	assert.Nil(t, err)

	devices, _ = db.GetDevices(profile.ID)

	assert.Equal(t, 3, len(devices))
	for _, d := range devices {
		assert.Equal(t, true, d.IsConfirmed)
	}

	db.RemoveDevice(profile.ID, "d1")
	db.RemoveDevice(profile.ID, "d2")

	devices, _ = db.GetDevices(profile.ID)
	assert.Equal(t, 1, len(devices))
	assert.Equal(t, "d3", devices[0].ID)

	db.RemoveDevice(profile.ID, "d3")
	devices, _ = db.GetDevices(profile.ID)
	assert.Equal(t, 0, len(devices))
}
