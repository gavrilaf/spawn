package dbx

import (
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateBridge(t *testing.T) {
	db, err := Connect(env.GetEnvironment("Test"))

	require.Nil(t, err)
	require.NotNil(t, db)
	require.NotNil(t, db.Db)

	db.Close()
}
