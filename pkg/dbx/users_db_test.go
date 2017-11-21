package dbx

import (
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	mdl "github.com/gavrilaf/spawn/pkg/dbx/model"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GetBridge(t *testing.T) Database {
	db, err := Connect(env.GetEnvironment("Test"))
	require.NotNil(t, db)
	require.Nil(t, err)

	return db
}

func TestRegisterProfile(t *testing.T) {
	db := GetBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	device := mdl.DeviceInfo{ID: "d1", Name: "d1-name", IsConfirmed: true}

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

	devices, err := db.GetUserDevices(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, devices)

	assert.Equal(t, 1, len(devices))
	assert.Equal(t, devices[0].ID, device.ID)
	assert.Equal(t, devices[0].UserID, profile.ID)
}

func TestReadAllProfiles(t *testing.T) {
	db := GetBridge(t)
	defer db.Close()
	profiles, errs := db.ReadAllUserProfiles()

	counter := 0
readLoop:
	for {
		select {
		case p := <-profiles:
			if p == nil {
				break readLoop
			}
			counter++
			continue readLoop
		case e := <-errs:
			assert.Truef(t, false, "Read profiles error: %v", e)
			break readLoop
		case <-time.After(10 * time.Second):
			assert.Truef(t, false, "Timeout")
			break readLoop
		}
	}

	fmt.Printf("Readed %d profiles\n", counter)
}

func TestDeviceManagement(t *testing.T) {
	db := GetBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	d1 := mdl.DeviceInfo{ID: "d1", Name: "d1-name", IsConfirmed: true}
	d2 := mdl.DeviceInfo{ID: "d2", Name: "d2-name", IsConfirmed: false}
	d3 := mdl.DeviceInfo{ID: "d3", Name: "d3-name", IsConfirmed: false}

	profile, err := db.RegisterUser(username, "password", d1)
	require.Nil(t, err)

	devices, err := db.GetUserDevices(profile.ID)
	require.Nil(t, err)

	assert.Equal(t, 1, len(devices))
	assert.Equal(t, devices[0].IsConfirmed, true)

	err = db.AddDevice(profile.ID, d2)
	require.Nil(t, err)

	err = db.AddDevice(profile.ID, d3)
	require.Nil(t, err)

	devices, _ = db.GetUserDevices(profile.ID)

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

	devices, _ = db.GetUserDevices(profile.ID)

	assert.Equal(t, 3, len(devices))
	for _, d := range devices {
		assert.Equal(t, true, d.IsConfirmed)
	}

	db.RemoveDevice(profile.ID, "d1")
	db.RemoveDevice(profile.ID, "d2")

	devices, _ = db.GetUserDevices(profile.ID)
	assert.Equal(t, 1, len(devices))
	assert.Equal(t, "d3", devices[0].ID)

	db.RemoveDevice(profile.ID, "d3")
	devices, _ = db.GetUserDevices(profile.ID)
	assert.Equal(t, 0, len(devices))
}

func TestEditUserProfile(t *testing.T) {
	db := GetBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	device := mdl.DeviceInfo{ID: "devive", Name: "d-name", IsConfirmed: true}

	profile, err := db.RegisterUser(username, "password", device)
	require.Nil(t, err)

	fmt.Printf("Before ****\n%v\n", spew.Sdump(profile))

	assert.Empty(t, profile.Country)

	permissions := profile.Permissions
	assert.False(t, permissions.Is2FARequired)
	assert.False(t, permissions.IsLocked)
	assert.False(t, permissions.IsEmailConfirmed)
	assert.Equal(t, int64(0), permissions.Scopes)

	personal := profile.PersonalInfo
	assert.Empty(t, personal.FirstName)
	assert.Empty(t, personal.LastName)
	assert.Equal(t, mdl.EmptyBirthDate, profile.GetBirthdayDate())

	err = db.UpdateUserCountry(profile.ID, "en")
	assert.Nil(t, err)

	permissions.Is2FARequired = true
	permissions.IsLocked = true
	permissions.IsEmailConfirmed = true
	permissions.Scopes = 2

	err = db.UpdateUserPermissions(profile.ID, permissions)
	assert.Nil(t, err)

	personal.FirstName = "FirstName"
	personal.LastName = "LastName"
	personal.BirthDate = mdl.BirthdayDate(1961, 10, 2)

	err = db.UpdateUserPersonalInfo(profile.ID, personal)
	assert.Nil(t, err)

	pr1, err := db.GetUserProfile(profile.ID)
	require.Nil(t, err)

	fmt.Printf("After ****\n%v\n", spew.Sdump(pr1))

	assert.Equal(t, "en", pr1.Country)

	assert.True(t, pr1.Is2FARequired)
	assert.True(t, pr1.IsLocked)
	assert.True(t, pr1.IsEmailConfirmed)
	assert.Equal(t, int64(2), pr1.Scopes)

	assert.Equal(t, "FirstName", pr1.FirstName)
	assert.Equal(t, "LastName", pr1.LastName)
	assert.Equal(t, mdl.BirthdayDate(1961, 10, 2), pr1.GetBirthdayDate())
}

func TestLoginLog(t *testing.T) {
	db := GetBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	device := mdl.DeviceInfo{ID: "device", Name: "d-name", IsConfirmed: true}

	profile, err := db.RegisterUser(username, "password", device)
	require.Nil(t, err)

	dex, err := db.GetUserDevicesEx(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, dex)

	fmt.Printf("Devices ex\n%v\n", spew.Sdump(dex))

	assert.Equal(t, 1, len(dex))
	assert.Equal(t, profile.ID, dex[0].UserID)
	assert.Equal(t, device.ID, dex[0].ID)
	assert.Empty(t, dex[0].GetLoginIP())

	err = db.LogUserLogin(profile.ID, "device", "127.0.0.1", "")
	assert.Nil(t, err)

	time.Sleep(time.Second)

	err = db.LogUserLogin(profile.ID, "device", "127.0.0.1", "")
	assert.Nil(t, err)

	dex, err = db.GetUserDevicesEx(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, dex)

	fmt.Printf("Devices ex\n%v\n", spew.Sdump(dex))

	assert.Equal(t, 1, len(dex))
	assert.Equal(t, profile.ID, dex[0].UserID)
	assert.Equal(t, device.ID, dex[0].ID)
	assert.Equal(t, "127.0.0.1", dex[0].GetLoginIP())

}
