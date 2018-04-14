package dbx

import (
	"time"

	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gavrilaf/spawn/pkg/utils"
)

func TestBridge_GetClients(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	clients, err := db.GetClients()
	assert.Nil(t, err)
	assert.NotEmpty(t, clients)
}

func TestBridge_UserNotFound(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	newID := uuid.NewV4().String()

	check := func(e error) {
		scope, reason := errx.GetErrorReason(e)
		assert.Equal(t, Scope, scope)
		assert.Equal(t, errx.ReasonNotFound, reason)
	}

	p, err := db.GetUserProfile(newID)
	assert.Nil(t, p)
	require.NotNil(t, err)
	check(err)

	p, err = db.FindUserProfile(newID)
	assert.Nil(t, p)
	require.NotNil(t, err)
	check(err)

	pd, err := db.GetUserDevice(newID, newID)
	assert.Nil(t, pd)
	require.NotNil(t, err)
	check(err)
}

func TestBridge_RegisterProfile(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	username := uuid.NewV4().String()

	device := testDevice

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

	device.UserID = profile.ID
	assert.Equal(t, device, devices[0])
}

func TestBridge_ReadAllProfiles(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	createTestUser(t, db)

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
		case e := <-errs:
			assert.Truef(t, false, "Read profiles error: %v", e)
			break readLoop
		case <-time.After(10 * time.Second):
			assert.Truef(t, false, "Timeout")
			break readLoop
		}
	}

	assert.True(t, counter > 0)
	//fmt.Printf("Readed %d profiles\n", counter)
}

func TestBridge_DeviceManagement(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	profile := createTestUser(t, db)

	// Check update device
	d1, _ := db.GetUserDevice(profile.ID, testDevice.DeviceID)

	err := db.SetDeviceFingerprint(profile.ID, d1.DeviceID, []byte("fingerprint"))
	assert.Nil(t, err)

	d1.Name = "d1-new-name"
	d1.Lang = "it"
	d1.Locale = "en"

	err = db.UpdateDevice(*d1)
	assert.Nil(t, err)

	dd1, _ := db.GetUserDevice(profile.ID, testDevice.DeviceID)

	d1.Fingerprint = []byte("fingerprint")

	assert.Equal(t, d1, dd1)

	// Check add/remove devices
	d2 := mdl.DeviceInfo{DeviceID: "d2", UserID: profile.ID, Name: "d2-name", IsConfirmed: false}
	d3 := mdl.DeviceInfo{DeviceID: "d3", UserID: profile.ID, Name: "d3-name", IsConfirmed: false}

	err = db.AddDevice(d2)
	require.Nil(t, err)

	err = db.AddDevice(d3)
	require.Nil(t, err)

	devices, _ := db.GetUserDevices(profile.ID)

	assert.Equal(t, 3, len(devices))
	for _, d := range devices {
		if d.DeviceID == "d2" || d.DeviceID == "d3" {
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

	db.RemoveDevice(profile.ID, d1.DeviceID)
	db.RemoveDevice(profile.ID, "d2")

	devices, _ = db.GetUserDevices(profile.ID)
	assert.Equal(t, 1, len(devices))
	assert.Equal(t, "d3", devices[0].DeviceID)

	db.RemoveDevice(profile.ID, "d3")

	devices, _ = db.GetUserDevices(profile.ID)
	assert.Equal(t, 0, len(devices))
}

func TestBridge_EditUserProfile(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	profile := createTestUser(t, db)

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

	err := db.UpdateUserCountry(profile.ID, "en")
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

	profile := createTestUser(t, db)

	dex, err := db.GetUserDevicesEx(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, dex)

	assert.Equal(t, 1, len(dex))
	assert.Equal(t, profile.ID, dex[0].UserID)
	assert.Equal(t, testDevice.DeviceID, dex[0].DeviceID)
	assert.Empty(t, dex[0].GetLoginIP())

	err = db.LogUserLogin(profile.ID, testDevice.DeviceID, "user-agent", "127.0.0.1", "")
	assert.Nil(t, err)

	time.Sleep(time.Second)

	err = db.LogUserLogin(profile.ID, testDevice.DeviceID, "user-agent", "127.0.0.1", "")
	assert.Nil(t, err)

	dex, err = db.GetUserDevicesEx(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, dex)

	assert.Equal(t, 1, len(dex))
	assert.Equal(t, profile.ID, dex[0].UserID)
	assert.Equal(t, testDevice.DeviceID, dex[0].DeviceID)
	assert.Equal(t, "127.0.0.1", dex[0].GetLoginIP())
}
