package cache

import (
	//"fmt"

	"testing"

	//"github.com/davecgh/go-spew/spew"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	//"github.com/gavrilaf/spawn/pkg/env"
	//"github.com/gavrilaf/spawn/pkg/errx"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

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
