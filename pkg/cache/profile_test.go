package cache

import (
	//"fmt"

	"testing"

	//"github.com/davecgh/go-spew/spew"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	//"github.com/gavrilaf/spawn/pkg/env"
	//"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gavrilaf/spawn/pkg/utils"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
	"database/sql"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/lib/pq"
	"time"
)

func TestBridge_UserProfile(t *testing.T) {
	cache := getCache(t)
	defer cache.Close()

	id := uuid.NewV4().String()

	profile := db.UserProfile{
		ID:      id,
		Country: "ua",
		PhoneNumber: db.PhoneNumber{
			CountryCode: 38,
			Number:      "067876123",
			IsConfirmed: false,
		},
		AuthInfo: db.AuthInfo{
			Username:     id + "@test.com",
			PasswordHash: "password",
			Permissions: db.Permissions{
				IsLocked:         true,
				IsEmailConfirmed: true,
				Is2FARequired:    true}},
		PersonalInfo: db.PersonalInfo{
			FirstName: "FirstName",
			LastName:  "LastName",
			BirthDate: utils.CreateDate(1966, 4, 21)}}

	err := cache.SetUserProfile(profile)
	assert.Nil(t, err)

	pr1, err := cache.GetUserProfile(id)
	assert.Nil(t, err)
	assert.NotNil(t, pr1)

	// TODO: Add all fields
	assert.Equal(t, id+"@test.com", pr1.Username)
	assert.Equal(t, "FirstName", pr1.FirstName)
	assert.Equal(t, "067876123", pr1.PhoneNumber.Number)

	//fmt.Printf("DB profile: %v\n", spew.Sdump(profile))
	//fmt.Printf("Cache profile: %v\n", spew.Sdump(pr1))

	assert.Equal(t, profile.BirthDate, pr1.GetBirthDate())
}

func TestBridge_ProfileNotFound(t *testing.T) {
	cache := getCache(t)
	defer cache.Close()

	id := uuid.NewV4().String() + "-not-found"

	p, err := cache.GetUserProfile(id)
	assert.Nil(t, p)
	assert.NotNil(t, err)

	scope, reason := errx.GetErrorReason(err)
	assert.Equal(t, Scope, scope)
	assert.Equal(t, errx.ReasonNotFound, reason)
}

func TestBridge_SetUserDevicesInfo(t *testing.T) {
	cache := getCache(t)
	defer cache.Close()

	devices := []db.DeviceInfoEx{
		db.DeviceInfoEx{
			LoginTime:   pq.NullTime{Time: time.Now(), Valid: true},
			LoginIP:     sql.NullString{String: "255.255.1.1", Valid: true},
			UserAgent:   sql.NullString{String: "test-22", Valid: true},
			LoginRegion: sql.NullString{String: "USA", Valid: true},
			DeviceInfo: db.DeviceInfo{
				ID:          "d1",
				IsConfirmed: true,
				Locale:      "en",
				Lang:        "en",
			},
		},
		db.DeviceInfoEx{
			LoginTime:   pq.NullTime{Time: time.Now(), Valid: true},
			LoginIP:     sql.NullString{String: "255.255.1.1", Valid: true},
			UserAgent:   sql.NullString{String: "test", Valid: true},
			LoginRegion: sql.NullString{String: "USA", Valid: true},
			DeviceInfo: db.DeviceInfo{
				ID:          "d2",
				IsConfirmed: true,
				Locale:      "es",
				Lang:        "es",
			},
		},
	}

	userID := uuid.NewV4().String()

	err := cache.SetUserDevicesInfo(userID, devices)
	assert.Nil(t, err)

	d2, err := cache.GetUserDevicesInfo(userID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(d2))

	err = cache.SetUserDevicesInfo(userID, nil)
	assert.Nil(t, err)

	d2, err = cache.GetUserDevicesInfo(userID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(d2))
}
