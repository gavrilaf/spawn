package cache

import (
	"testing"

	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"

	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridge_AddClient(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	cl := db.Client{"cl-1", []byte("secret"), true, "desc", 0}

	err := cache.AddClient(cl)
	require.Nil(t, err)

	p, err := cache.FindClient(cl.ID)
	assert.Nil(t, err)
	assert.NotNil(t, p)

	assert.Equal(t, cl, *p)
}

func TestBridge_AddSession(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	session := mdl.Session{
		RefreshToken:      "refresh-token",
		ClientID:          "client-id",
		ClientSecret:      []byte("secret"),
		UserID:            "user-id-11",
		DeviceID:          "device-id-11",
		IsDeviceConfirmed: true,
		Locale:            "en",
		Lang:              "en"}

	sessionID, err := cache.AddSession(session, false)
	require.Nil(t, err)

	p, err := cache.GetSession(sessionID)
	require.Nil(t, err)
	require.NotNil(t, p)

	assert.Equal(t, sessionID, p.ID)
	assert.Equal(t, session.ClientID, p.ClientID)
	assert.Equal(t, session.RefreshToken, p.RefreshToken)
	assert.Equal(t, session.ClientSecret, p.ClientSecret)
	assert.Equal(t, session.UserID, p.UserID)
	assert.Equal(t, session.DeviceID, p.DeviceID)
	assert.Equal(t, session.IsDeviceConfirmed, p.IsDeviceConfirmed)
	assert.Equal(t, session.Locale, p.Locale)
	assert.Equal(t, session.Lang, p.Lang)

	err = cache.DeleteSession(sessionID)
	assert.Nil(t, err)

	p, err = cache.GetSession(sessionID)
	require.NotNil(t, err)
	require.Nil(t, p)
}

func TestBridge_AddSessionIfExists(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	session := mdl.Session{
		RefreshToken:      "refresh-token",
		ClientID:          "client-id",
		ClientSecret:      []byte("secret"),
		UserID:            "user-id-21",
		DeviceID:          "device-id-21",
		IsDeviceConfirmed: true,
		Locale:            "en",
		Lang:              "en"}

	sessionID, err := cache.AddSession(session, false)
	require.Nil(t, err)

	sessionID, err = cache.AddSession(session, false)
	assert.NotNil(t, err)

	scope, reason := errx.GetErrorReason(err)
	assert.Equal(t, Scope, scope)
	assert.Equal(t, ReasonSessionDuplicate, reason)

	sessionID2, err := cache.AddSession(session, true)
	assert.Nil(t, err)

	assert.NotEqual(t, sessionID, sessionID2)

	p, err := cache.GetSession(sessionID)
	require.NotNil(t, err)
	require.Nil(t, p)

	p, err = cache.GetSession(sessionID2)
	require.Nil(t, err)
	require.NotNil(t, p)

	err = cache.DeleteSession(sessionID2)
	assert.Nil(t, err)
}

func TestBridge_FindSession(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	session := mdl.Session{
		RefreshToken:      "refresh-token",
		ClientID:          "client-id",
		ClientSecret:      []byte("secret"),
		UserID:            "user-id-1111",
		DeviceID:          "device-id-1111",
		IsDeviceConfirmed: true,
		Locale:            "en",
		Lang:              "en"}

	sessionID, err := cache.AddSession(session, false)
	require.Nil(t, err)

	p, err := cache.FindSession("user-id-1111", "device-id-1111")
	assert.Nil(t, err)
	assert.NotNil(t, p)

	assert.Equal(t, sessionID, p.ID)
	assert.Equal(t, session.ClientID, p.ClientID)
	assert.Equal(t, session.RefreshToken, p.RefreshToken)
	assert.Equal(t, session.ClientSecret, p.ClientSecret)
	assert.Equal(t, session.UserID, p.UserID)
	assert.Equal(t, session.DeviceID, p.DeviceID)
	assert.Equal(t, session.IsDeviceConfirmed, p.IsDeviceConfirmed)
	assert.Equal(t, session.Locale, p.Locale)
	assert.Equal(t, session.Lang, p.Lang)

	err = cache.DeleteSession(sessionID)
	assert.Nil(t, err)
}

func TestBridge_SetUserAuthInfo(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	profile := db.UserProfile{
		ID: "user-1",
		AuthInfo: db.AuthInfo{
			Username:     "testuser@test.com",
			PasswordHash: "password",
			Permissions: db.Permissions{
				IsLocked:         true,
				IsEmailConfirmed: true,
				Is2FARequired:    true}},
		PersonalInfo: db.PersonalInfo{
			FirstName: "FirstName",
			LastName:  "LastName"}}

	devices := []db.DeviceInfo{
		db.DeviceInfo{ID: "d1"},
		db.DeviceInfo{ID: "id2", Fingerprint: []byte("fingerpring")},
	}

	err := cache.SetUserAuthInfo(profile, devices)
	require.Nil(t, err)

	p1, err := cache.FindUserAuthInfo(profile.Username)
	require.Nil(t, err)
	require.NotNil(t, p1)

	assert.Equal(t, profile.ID, p1.ID)
	assert.Equal(t, profile.Username, p1.Username)
	assert.Equal(t, profile.IsLocked, p1.IsLocked)
	assert.Equal(t, profile.IsEmailConfirmed, p1.IsEmailConfirmed)
	assert.Equal(t, profile.Is2FARequired, p1.Is2FARequired)

	p2, err := cache.FindUserAuthInfo("unknown-user-name@@@")
	require.NotNil(t, err)
	require.Nil(t, p2)
}

func TestBridge_SetDevice(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	profile := db.UserProfile{
		ID: "user-1",
		AuthInfo: db.AuthInfo{
			Username:     "testuser@test.com",
			PasswordHash: "password",
			Permissions: db.Permissions{
				IsLocked:         false,
				IsEmailConfirmed: false,
				Is2FARequired:    false}}}

	devices := []db.DeviceInfo{
		db.DeviceInfo{ID: "d1", IsConfirmed: false, Locale: "ru", Lang: "ru"},
		db.DeviceInfo{ID: "d2", IsConfirmed: true, Fingerprint: []byte("fingerprint"), Locale: "en", Lang: "en"},
	}

	err := cache.SetUserAuthInfo(profile, devices)
	require.Nil(t, err)

	d1, err := cache.GetDevice("user-1", "d1")
	assert.Nil(t, err)
	assert.NotNil(t, d1)

	d2, err := cache.GetDevice("user-1", "d2")
	assert.Nil(t, err)

	assert.Equal(t, "d2", d2.DeviceID)
	assert.Equal(t, "user-1", d2.UserID)
	assert.Equal(t, true, d2.IsConfirmed)
	assert.Equal(t, []byte("fingerprint"), d2.Fingerprint)
	assert.Equal(t, "en", d2.Locale)
	assert.Equal(t, "en", d2.Lang)

	assert.Nil(t, cache.DeleteDevice(profile.ID, "d1"))
	dd1, _ := cache.GetDevice(profile.ID, "d1")
	assert.Nil(t, dd1)

	err = cache.SetDevice(db.DeviceInfo{ID: "d3", UserID: profile.ID})
	assert.Nil(t, err)

	d3, _ := cache.GetDevice(profile.ID, "d3")
	assert.NotNil(t, d3)

	assert.Nil(t, cache.DeleteDevice(profile.ID, "d3"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d1"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d2"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d2"))
}

func TestBridge_AuthNotFound(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	new_id := uuid.NewV4().String() + "-not-found"

	check := func(e error) {
		scope, reason := errx.GetErrorReason(e)
		assert.Equal(t, Scope, scope)
		assert.Equal(t, errx.ReasonNotFound, reason)
	}

	p, err := cache.FindClient(new_id)
	assert.Nil(t, p)
	require.NotNil(t, err)
	check(err)

	p1, err := cache.GetSession(new_id)
	assert.Nil(t, p1)
	require.NotNil(t, err)
	check(err)

	p2, err := cache.FindUserAuthInfo(new_id)
	assert.Nil(t, p2)
	require.NotNil(t, err)
	check(err)

	p3, err := cache.GetDevice(new_id, new_id)
	assert.Nil(t, p3)
	require.NotNil(t, err)
	check(err)

	p4, err := cache.GetUserProfile(new_id)
	assert.Nil(t, p4)
	require.NotNil(t, err)
	check(err)
}
