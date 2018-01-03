package cache

import (
	"testing"
	"time"

	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridge_SetUserAccounts(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	userID := uuid.NewV4().String()

	tm := time.Now().UTC().Truncate(time.Second)

	accounts := []db.Account{
		{ID: uuid.NewV4().String(), UserID: userID, Name: "BTC wallet", Currency: "BTC", Created: tm},
		{ID: uuid.NewV4().String(), UserID: userID, Name: "USD wallet", Currency: "USD", Created: tm},
	}

	err := cache.SetUserAccounts(userID, accounts)
	assert.Nil(t, err)

	accounts2, err := cache.GetUserAccounts(userID)
	assert.Nil(t, err)
	require.Equal(t, 2, len(accounts2))

	assert.NotEmpty(t, accounts2[0].ID)
	assert.NotEmpty(t, accounts2[0].Name)
	assert.NotEmpty(t, accounts2[0].Currency)

	assert.Equal(t, tm, accounts2[0].GetCreated())

	assert.Equal(t, mdl.ACCOUNT_CREATED, accounts2[0].AccountState.Status)
	assert.Equal(t, int64(0), accounts2[0].AccountState.Updated)
	assert.Empty(t, accounts2[0].AccountState.Balance)

	err = cache.ClearUserAccounts(userID)
	assert.Nil(t, err)

	accounts2, err = cache.GetUserAccounts(userID)
	assert.Nil(t, err)
	require.Equal(t, 0, len(accounts2))
}

func TestBridge_UpdateAccount(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	userID := uuid.NewV4().String()

	tm := time.Now().UTC().Truncate(time.Second)

	accounts := []db.Account{
		{ID: uuid.NewV4().String(), UserID: userID, Name: "BTC wallet", Currency: "BTC", Created: tm},
	}

	err := cache.SetUserAccounts(userID, accounts)
	assert.Nil(t, err)

	account, err := cache.GetUserAccount(userID, accounts[0].ID)
	assert.Nil(t, err)
	require.NotNil(t, account)

	assert.Equal(t, accounts[0].ID, account.ID)
	assert.Equal(t, accounts[0].Name, account.Name)
	assert.Equal(t, accounts[0].Currency, account.Currency)
	assert.Equal(t, accounts[0].Created.UTC().Unix(), account.Created)

	assert.Equal(t, mdl.ACCOUNT_CREATED, account.AccountState.Status)
	assert.Equal(t, int64(0), account.AccountState.Updated)
	assert.Empty(t, account.AccountState.Balance)

	err = cache.UpdateUserAccountBalance(userID, account.ID, "0.011")
	assert.Nil(t, err)

	err = cache.UpdateUserAccountStatus(userID, account.ID, mdl.ACCOUNT_ACTIVE)
	assert.Nil(t, err)

	account, err = cache.GetUserAccount(userID, accounts[0].ID)
	assert.Nil(t, err)

	assert.Equal(t, mdl.ACCOUNT_ACTIVE, account.AccountState.Status)
	assert.Equal(t, "0.011", account.AccountState.Balance)
	assert.NotEqual(t, int64(0), account.AccountState.Updated)

	err = cache.ClearUserAccounts(userID)
	assert.Nil(t, err)
}

func TestBridge_AddAccount(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	userID := uuid.NewV4().String()

	tm := time.Now().UTC().Truncate(time.Second)

	accounts := []db.Account{
		{ID: uuid.NewV4().String(), UserID: userID, Name: "BTC wallet", Currency: "BTC", Created: tm},
	}

	err := cache.SetUserAccounts(userID, accounts)
	assert.Nil(t, err)

	accounts2, err := cache.GetUserAccounts(userID)
	require.Equal(t, 1, len(accounts2))

	new_id := uuid.NewV4().String()
	new_acc := db.Account{ID: new_id, UserID: userID, Name: "USD wallet", Currency: "USD", Created: tm}

	err = cache.AddUserAccount(userID, new_acc)
	assert.Nil(t, err)

	accounts2, err = cache.GetUserAccounts(userID)
	require.Equal(t, 2, len(accounts2))

	found := false
	for _, p := range accounts2 {
		if p.ID == new_id {
			found = true

			assert.Equal(t, "USD", p.Currency)
			assert.Equal(t, "USD wallet", p.Name)
			assert.Equal(t, tm, p.GetCreated())

			assert.Equal(t, mdl.ACCOUNT_CREATED, p.AccountState.Status)
			assert.Equal(t, int64(0), p.AccountState.Updated)
			assert.Empty(t, p.AccountState.Balance)
		}
	}

	assert.True(t, found)

	err = cache.ClearUserAccounts(userID)
	assert.Nil(t, err)
}

func TestBridge_AccountNotFound(t *testing.T) {
	cache := getTestCache(t)
	defer cache.Close()

	user_id := uuid.NewV4().String()
	account_id := uuid.NewV4().String()

	accounts, err := cache.GetUserAccounts(user_id)
	assert.Nil(t, err)
	assert.Empty(t, accounts)

	p, err := cache.GetUserAccount(user_id, account_id)
	assert.NotNil(t, err)
	assert.Nil(t, p)

	scope, reason := errx.GetErrorReason(err)
	assert.Equal(t, Scope, scope)
	assert.Equal(t, errx.ReasonNotFound, reason)
}
