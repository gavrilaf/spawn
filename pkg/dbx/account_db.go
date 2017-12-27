package dbx

import (
	"time"

	"github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/satori/go.uuid"
)

const (
	getSupportedAccounts = `SELECT * FROM public."AccountMeta"`
	addAccount           = `INSERT INTO public."Accounts"(id, user_id, name, currency, created) VALUES ($1, $2, $3, $4, $5)`
	getAccountByID       = `SELECT * FROM public."Accounts" WHERE id = $1`
	getUserAccounts      = `SELECT * FROM public."Accounts" WHERE user_id = $1`
)

func (db *Bridge) GetSupportedAccounts() ([]mdl.AccountMeta, error) {
	var accountMeta []mdl.AccountMeta
	if err := db.conn.Select(&accountMeta, getSupportedAccounts); err != nil {
		return nil, err
	}
	return accountMeta, nil
}

func (db *Bridge) GetAccountsAllowedForUser(userID string) ([]mdl.AccountMeta, error) {
	supported, err := db.GetSupportedAccounts()
	if err != nil {
		return nil, err
	}

	userAccounts, err := db.GetUserAccounts(userID)
	if err != nil {
		return nil, err
	}

	isRegistered := func(curr string) bool {
		for _, acc := range userAccounts {
			if acc.Currency == curr {
				return true
			}
		}
		return false
	}

	var ret []mdl.AccountMeta
	for _, meta := range supported {
		if meta.MultipleAllowed || !isRegistered(meta.Currency) {
			ret = append(ret, meta)
		}
	}

	return ret, nil
}

func (db *Bridge) GetUserAccounts(userID string) ([]mdl.Account, error) {
	var accounts []mdl.Account
	if err := db.conn.Select(&accounts, getUserAccounts, userID); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (db *Bridge) AddAccount(userID string, currency string, name string) (*mdl.Account, error) {
	accountID := uuid.NewV4().String()
	_, err := db.conn.Exec(addAccount, accountID, userID, name, currency, time.Now().UTC())
	if err != nil {

		return nil, err
	}

	return db.GetAccount(accountID)
}

func (db *Bridge) GetAccount(accountID string) (*mdl.Account, error) {
	var account mdl.Account
	if err := db.conn.Get(&account, getAccountByID, accountID); err != nil {
		return nil, errx.ErrNotFound(Scope, "Account with id %v not found: %v", accountID, err)
	}

	return &account, nil
}
