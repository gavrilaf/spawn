package cache

import (
	//"fmt"
	"testing"

	mdl "github.com/gavrilaf/spawn/pkg/dbx/model"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GetEnv() *env.Environment {
	return env.GetEnvironment("Test")
}

func TestClientCache(t *testing.T) {
	cache, err := Connect(GetEnv())
	require.Nil(t, err)
	require.NotNil(t, cache)

	defer cache.Close()

	cl := mdl.Client{"cl-1", []byte("secret")}

	err = cache.AddClient(cl)
	require.Nil(t, err)

	p, err := cache.FindClient(cl.ID)
	assert.Nil(t, err)
	assert.NotNil(t, p)

	assert.Equal(t, cl.ID, p.ID)
	assert.Equal(t, cl.Secret, p.Secret)

	p, err = cache.FindClient("unexisting-client-id-")
	assert.NotNil(t, err)
	assert.Nil(t, p)
}

func TestUserSession(t *testing.T) {
	cache, err := Connect(GetEnv())
	require.Nil(t, err)
	require.NotNil(t, cache)

	defer cache.Close()

	session := Session{
		ID:                "ses-1",
		RefreshToken:      "refresh-token",
		ClientID:          "client-id",
		ClientSecret:      []byte("secret"),
		UserID:            "user-id",
		DeviceID:          "device-id",
		IsDeviceConfirmed: true,
		Locale:            "en",
		Lang:              "en"}

	err = cache.AddSession(session)
	require.Nil(t, err)

	p, err := cache.FindSession(session.ID)
	require.Nil(t, err)
	require.NotNil(t, p)

	assert.Equal(t, session.ID, p.ID)
	assert.Equal(t, session.ClientID, p.ClientID)
	assert.Equal(t, session.RefreshToken, p.RefreshToken)
	assert.Equal(t, session.ClientSecret, p.ClientSecret)
	assert.Equal(t, session.UserID, p.UserID)
	assert.Equal(t, session.DeviceID, p.DeviceID)
	assert.Equal(t, session.IsDeviceConfirmed, p.IsDeviceConfirmed)
	assert.Equal(t, session.Locale, p.Locale)
	assert.Equal(t, session.Lang, p.Lang)

	err = cache.DeleteSession(session.ID)
	assert.Nil(t, err)

	p, err = cache.FindSession(session.ID)
	//fmt.Printf("Error: %v\n", err)
	require.NotNil(t, err)
	require.Nil(t, p)
}

func TestAuthUser(t *testing.T) {
	cache, err := Connect(GetEnv())
	require.Nil(t, err)
	require.NotNil(t, cache)

	defer cache.Close()

	profile := mdl.UserProfile{
		ID: "user-1",
		AuthInfo: mdl.AuthInfo{
			Username:     "testuser@test.com",
			PasswordHash: "password",
			Permissions: mdl.Permissions{
				IsLocked:         true,
				IsEmailConfirmed: true,
				Is2FARequired:    true}},
		PersonalInfo: mdl.PersonalInfo{
			FirstName: "FirstName",
			LastName:  "LastName"}}

	devices := []mdl.DeviceInfo{
		mdl.DeviceInfo{ID: "d1"},
		mdl.DeviceInfo{ID: "id2", Fingerprint: []byte("fingerpring")},
	}

	err = cache.SetUserAuthInfo(profile, devices)
	require.Nil(t, err)

	p1, err := cache.FindUserAuthInfo(profile.Username)
	require.Nil(t, err)
	require.NotNil(t, p1)

	assert.Equal(t, profile.ID, p1.ID)
	assert.Equal(t, profile.Username, p1.Username)
	assert.Equal(t, profile.IsLocked, p1.IsLocked)
	assert.Equal(t, profile.IsEmailConfirmed, p1.IsEmailConfirmed)
	assert.Equal(t, profile.Is2FARequired, p1.Is2FARequired)
}

func TestUserDevices(t *testing.T) {
	cache, err := Connect(GetEnv())
	require.Nil(t, err)
	defer cache.Close()

	profile := mdl.UserProfile{
		ID: "user-1",
		AuthInfo: mdl.AuthInfo{
			Username:     "testuser@test.com",
			PasswordHash: "password",
			Permissions: mdl.Permissions{
				IsLocked:         false,
				IsEmailConfirmed: false,
				Is2FARequired:    false}},
		PersonalInfo: mdl.PersonalInfo{
			FirstName: "FirstName",
			LastName:  "LastName"}}

	devices := []mdl.DeviceInfo{
		mdl.DeviceInfo{ID: "d1", IsConfirmed: false, Locale: "ru", Lang: "ru"},
		mdl.DeviceInfo{ID: "d2", IsConfirmed: true, Fingerprint: []byte("fingerprint"), Locale: "en", Lang: "en"},
	}

	err = cache.SetUserAuthInfo(profile, devices)
	require.Nil(t, err)

	d1, err := cache.FindDevice("user-1", "d1")
	assert.Nil(t, err)
	assert.NotNil(t, d1)

	d2, err := cache.FindDevice("user-1", "d2")
	assert.Nil(t, err)

	assert.Equal(t, "d2", d2.DeviceID)
	assert.Equal(t, "user-1", d2.UserID)
	assert.Equal(t, true, d2.IsConfirmed)
	assert.Equal(t, []byte("fingerprint"), d2.Fingerpring)
	assert.Equal(t, "en", d2.Locale)
	assert.Equal(t, "en", d2.Lang)

	assert.Nil(t, cache.DeleteDevice(profile.ID, "d1"))
	dd1, _ := cache.FindDevice(profile.ID, "d1")
	assert.Nil(t, dd1)

	err = cache.SetDevice(profile.ID, mdl.DeviceInfo{ID: "d3"})
	assert.Nil(t, err)

	d3, _ := cache.FindDevice(profile.ID, "d3")
	assert.NotNil(t, d3)

	assert.Nil(t, cache.DeleteDevice(profile.ID, "d3"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d1"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d2"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d2"))
}

func TestConfirmCodes(t *testing.T) {
	cache, err := Connect(GetEnv())
	require.Nil(t, err)
	defer cache.Close()

	err = cache.AddConfirmCode("device", "d-id-1", "123456")
	assert.Nil(t, err)

	code, err := cache.GetConfirmCode("device", "d-id-1")
	assert.Nil(t, err)
	assert.Equal(t, "123456", code)

	err = cache.DeleteConfirmCode("device", "d-id-1")
	assert.Nil(t, err)

	code, _ = cache.GetConfirmCode("device", "d-id-1")
	assert.Equal(t, "", code)
}
