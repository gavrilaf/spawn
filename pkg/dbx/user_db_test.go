package dbx

import (
	"fmt"
	"testing"
	"time"

	"github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gavrilaf/spawn/pkg/utils"

	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridge_GetClients(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	clients, err := db.GetClients()
	assert.Nil(t, err)
	assert.NotEmpty(t, clients)
}

func TestBridge_ProfileNotFound(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	new_id := uuid.NewV4().String()

	check := func(e error) {
		scope, reason := errx.GetErrorReason(e)
		assert.Equal(t, Scope, scope)
		assert.Equal(t, errx.ReasonNotFound, reason)
	}

	p, err := db.GetUserProfile(new_id)
	assert.Nil(t, p)
	require.NotNil(t, err)
	check(err)

	p, err = db.FindUserProfile(new_id)
	assert.Nil(t, p)
	require.NotNil(t, err)
	check(err)

	pd, err := db.GetUserDevice(new_id, new_id)
	assert.Nil(t, pd)
	require.NotNil(t, err)
	check(err)
}

func TestBridge_RegisterProfile(t *testing.T) {
	db := getBridge(t)
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

func TestBridge_ReadAllProfiles(t *testing.T) {
	db := getBridge(t)
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

func TestBridge_DeviceManagement(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	d1 := mdl.DeviceInfo{ID: "d1", Name: "d1-name", IsConfirmed: true, Lang: "ru", Locale: "es"}

	profile, err := db.RegisterUser(username, "password", d1)
	require.Nil(t, err)

	d2 := mdl.DeviceInfo{ID: "d2", UserID: profile.ID, Name: "d2-name", IsConfirmed: false}
	d3 := mdl.DeviceInfo{ID: "d3", UserID: profile.ID, Name: "d3-name", IsConfirmed: false}

	devices, err := db.GetUserDevices(profile.ID)
	require.Nil(t, err)

	assert.Equal(t, 1, len(devices))
	assert.Equal(t, true, devices[0].IsConfirmed)
	assert.Equal(t, "d1-name", devices[0].Name)
	assert.Equal(t, "ru", devices[0].Lang)
	assert.Equal(t, "es", devices[0].Locale)

	d1.UserID = profile.ID
	d1.Name = "d1-new-name"
	d1.Lang = "it"
	d1.Locale = "en"

	assert.Equal(t, []byte(nil), d1.Fingerprint)

	err = db.SetDeviceFingerprint(profile.ID, "d1", []byte("fingerprint"))
	assert.Nil(t, err)

	err = db.UpdateDevice(d1)
	assert.Nil(t, err)

	dd1, err := db.GetUserDevice(profile.ID, "d1")
	assert.Nil(t, err)
	require.NotNil(t, dd1)

	assert.Equal(t, true, dd1.IsConfirmed)
	assert.Equal(t, "d1-new-name", dd1.Name)
	assert.Equal(t, "it", dd1.Lang)
	assert.Equal(t, "en", dd1.Locale)

	assert.Equal(t, []byte("fingerprint"), dd1.Fingerprint)

	err = db.AddDevice(d2)
	require.Nil(t, err)

	err = db.AddDevice(d3)
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

func TestBridge_EditUserProfile(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	device := mdl.DeviceInfo{ID: "devive", Name: "d-name", IsConfirmed: true}

	profile, err := db.RegisterUser(username, "password", device)
	require.Nil(t, err)

	assert.Empty(t, profile.Country)

	phoneNumber := profile.PhoneNumber
	assert.Equal(t, 0, phoneNumber.CountryCode)
	assert.Equal(t, "", phoneNumber.Number)
	assert.Equal(t, false, phoneNumber.IsConfirmed)

	permissions := profile.Permissions
	assert.False(t, permissions.Is2FARequired)
	assert.False(t, permissions.IsLocked)
	assert.False(t, permissions.IsEmailConfirmed)
	assert.Equal(t, 0, permissions.Scope)

	personal := profile.PersonalInfo
	assert.Empty(t, personal.FirstName)
	assert.Empty(t, personal.LastName)
	assert.Equal(t, utils.EmptyBirthdayDate, personal.BirthDate)

	err = db.UpdateUserCountry(profile.ID, "en")
	assert.Nil(t, err)

	permissions.Is2FARequired = true
	permissions.IsLocked = true
	permissions.IsEmailConfirmed = true
	permissions.Scope = 2

	err = db.UpdateUserPermissions(profile.ID, permissions)
	assert.Nil(t, err)

	personal.FirstName = "FirstName"
	personal.LastName = "LastName"
	personal.BirthDate = utils.CreateDate(1961, 10, 2)

	err = db.UpdateUserPersonalInfo(profile.ID, personal)
	assert.Nil(t, err)

	phoneNumber.CountryCode = 38
	phoneNumber.Number = "97822345"

	err = db.UpdateUserPhoneNumber(profile.ID, phoneNumber)
	assert.Nil(t, err)

	pr1, err := db.GetUserProfile(profile.ID)
	require.Nil(t, err)

	assert.Equal(t, "en", pr1.Country)

	assert.True(t, pr1.Is2FARequired)
	assert.True(t, pr1.IsLocked)
	assert.True(t, pr1.IsEmailConfirmed)
	assert.Equal(t, 2, pr1.Scope)

	assert.Equal(t, "FirstName", pr1.FirstName)
	assert.Equal(t, "LastName", pr1.LastName)
	assert.Equal(t, utils.CreateDate(1961, 10, 2), pr1.BirthDate)

	phoneNumber = pr1.PhoneNumber
	assert.Equal(t, 38, phoneNumber.CountryCode)
	assert.Equal(t, "97822345", phoneNumber.Number)
}

func TestBridge_LoginLog(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	device := mdl.DeviceInfo{ID: "device", Name: "d-name", IsConfirmed: true}

	profile, err := db.RegisterUser(username, "password", device)
	require.Nil(t, err)

	dex, err := db.GetUserDevicesEx(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, dex)

	assert.Equal(t, 1, len(dex))
	assert.Equal(t, profile.ID, dex[0].UserID)
	assert.Equal(t, device.ID, dex[0].ID)
	assert.Empty(t, dex[0].GetLoginIP())

	err = db.LogUserLogin(profile.ID, "device", "user-agent", "127.0.0.1", "")
	assert.Nil(t, err)

	time.Sleep(time.Second)

	err = db.LogUserLogin(profile.ID, "device", "user-agent", "127.0.0.1", "")
	assert.Nil(t, err)

	dex, err = db.GetUserDevicesEx(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, dex)

	assert.Equal(t, 1, len(dex))
	assert.Equal(t, profile.ID, dex[0].UserID)
	assert.Equal(t, device.ID, dex[0].ID)
	assert.Equal(t, "127.0.0.1", dex[0].GetLoginIP())
}
