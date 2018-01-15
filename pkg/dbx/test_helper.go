package dbx

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func getBridge(t *testing.T) Database {
	db, err := Connect(senv.GetEnvironment("Test"))
	require.NotNil(t, db)
	require.Nil(t, err)

	return db
}

func deleteAccount(db Database, accountID string) error {
	conn := db.(*Bridge).conn
	_, err := conn.Exec(`DELETE FROM public."Accounts" WHERE ID = $1`, accountID)
	return err
}

func createTestUser(db Database) (*mdl.UserProfile, error) {
	username := uuid.NewV4().String()
	device := mdl.DeviceInfo{ID: "d1-tes", Name: "d1-test-name", IsConfirmed: true}
	return db.RegisterUser(username, "password", device)
}
