package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"fmt"
	"github.com/gavrilaf/spawn/pkg/env"
)

func GetEnv() *env.Environment {
	return env.GetEnvironment("Test")
}

func TestClientCache(t *testing.T) {
	cache, err := Connect(GetEnv())
	require.Nil(t, err)
	require.NotNil(t, cache)

	defer cache.Close()

	cl := Client{"cl-1", []byte("secret")}

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

	session := Session{"ses-1", "refresh-token", "client-id", []byte("secret"), "user-id", "device-id"}

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

	err = cache.DeleteSession(session.ID)
	assert.Nil(t, err)

	p, err = cache.FindSession(session.ID)
	fmt.Printf("Error: %v\n", err)
	require.NotNil(t, err)
	require.Nil(t, p)
}

func TestUserProfile(t *testing.T) {
	cache, err := Connect(GetEnv())
	require.Nil(t, err)
	require.NotNil(t, cache)

	defer cache.Close()

	profile := UserProfile{"user-1", AuthInfo{"testuser@test.com", "password", false}, PersonalInfo{"FirstName", "LastName"}}
	devices := []string{"device-1", "device-2"}

	err = cache.AddUser(profile, devices)
	require.Nil(t, err)

	p1, err := cache.FindProfile(profile.ID)
	require.Nil(t, err)
	require.NotNil(t, p1)

	assert.Equal(t, profile.ID, p1.ID)
	assert.Equal(t, profile.Username, p1.Username)
	assert.Equal(t, profile.IsLocked, p1.IsLocked)
	assert.Equal(t, profile.LastName, p1.LastName)
}

func TestUserDevices(t *testing.T) {
	cache, err := Connect(GetEnv())
	require.Nil(t, err)
	require.NotNil(t, cache)

	defer cache.Close()

	profile := UserProfile{"user-2", AuthInfo{"testuser@test.com", "password", false}, PersonalInfo{"FirstName", "LastName"}}
	devices := []string{"d1", "d2"}

	err = cache.AddUser(profile, devices)
	require.Nil(t, err)

	b, err := cache.IsDeviceExists(profile.ID, "d1")
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	assert.Nil(t, cache.DeleteDevice(profile.ID, "d1"))
	b, _ = cache.IsDeviceExists(profile.ID, "d1")
	assert.Equal(t, false, b)

	b, _ = cache.IsDeviceExists(profile.ID, "d3")
	assert.Equal(t, false, b)

	err = cache.AddDevice(profile.ID, "d3")
	assert.Nil(t, err)

	b, _ = cache.IsDeviceExists(profile.ID, "d3")
	assert.Equal(t, true, b)

	assert.Nil(t, cache.DeleteDevice(profile.ID, "d3"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d1"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d2"))
	assert.Nil(t, cache.DeleteDevice(profile.ID, "d2"))
}
