package cache

import (
	"database/sql"
	"testing"
	"time"

	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestBridge_ConfirmCode(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	err := cache.AddConfirmCode("device", "d-id-1", "123456")
	assert.Nil(t, err)

	code, err := cache.GetConfirmCode("device", "d-id-1")
	assert.Nil(t, err)
	assert.Equal(t, "123456", code)

	err = cache.DeleteConfirmCode("device", "d-id-1")
	assert.Nil(t, err)

	code, _ = cache.GetConfirmCode("device", "d-id-1")
	assert.Equal(t, "", code)
}

func TestBridge_SetUserDevicesInfo(t *testing.T) {
	cache := getTestCache(t)
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
