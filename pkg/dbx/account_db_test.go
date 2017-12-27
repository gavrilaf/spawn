package dbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridge_GetSupportedAccounts(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	supported, err := db.GetSupportedAccounts()
	assert.Nil(t, err)
	assert.NotEmpty(t, supported)

	assert.Equal(t, 3, len(supported))
}

func TestBridge_AddAccount(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	user, _ := createTestUser(db)
	require.NotNil(t, user)

	acc1, err := db.AddAccount(user.ID, "BTC", "BTC test wallet")
	assert.Nil(t, err)
	require.NotNil(t, acc1)

	acc2, err := db.GetAccount(acc1.ID)
	assert.Nil(t, err)
	require.NotNil(t, acc2)

	assert.Equal(t, user.ID, acc1.UserID)
	assert.Equal(t, "BTC", acc1.Currency)
	assert.Equal(t, "BTC test wallet", acc1.Name)

	assert.Equal(t, acc1.ID, acc2.ID)
	assert.Equal(t, acc1.UserID, acc2.UserID)
	assert.Equal(t, acc1.Currency, acc2.Currency)
	assert.Equal(t, acc1.Name, acc2.Name)
}

func TestBridge_GetUserAccounts(t *testing.T) {
	db := getBridge(t)
	defer db.Close()

	user, _ := createTestUser(db)
	require.NotNil(t, user)

	allowed, err := db.GetAccountsAllowedForUser(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(allowed))

	_, err = db.AddAccount(user.ID, "USD", "USD test wallet")
	assert.Nil(t, err)

	accounts, err := db.GetUserAccounts(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(accounts))

	assert.Equal(t, user.ID, accounts[0].UserID)
	assert.Equal(t, "USD", accounts[0].Currency)
	assert.Equal(t, "USD test wallet", accounts[0].Name)

	allowed, err = db.GetAccountsAllowedForUser(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(allowed))

	deleteAccount(db, accounts[0].ID)

	accounts, err = db.GetUserAccounts(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(accounts))

	allowed, err = db.GetAccountsAllowedForUser(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(allowed))
}
