package dbx

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func getBridge(t *testing.T) Database {
	db, err := Connect(senv.GetEnvironment())
	require.NotNil(t, db)
	require.Nil(t, err)

	return db
}

func deleteAccount(db Database, accountID string) error {
	conn := db.(*Bridge).conn
	_, err := conn.Exec(`DELETE FROM public."Accounts" WHERE ID = $1`, accountID)
	return err
}

var testDevice = mdl.DeviceInfo{DeviceID: "d1-tes", Name: "d1-test-name", IsConfirmed: true, Fingerprint: []uint8{}, Locale: "ua", Lang: "en"}

func createTestUser(t *testing.T, db Database) *mdl.UserProfile {
	username := uuid.NewV4().String()
	p, err := db.RegisterUser(username, "password", testDevice)
	require.Nil(t, err)

	return p
}
