package dbx

import (
	"github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	getSupportedAccounts = ""
)

func (db *Bridge) GetSupportedAccounts() ([]mdl.AccountMeta, error) {
	var accountMetas []mdl.AccountMeta
	if err := db.conn.Select(&accountMetas, getSupportedAccounts); err != nil {
		return nil, err
	}
	return accountMetas, nil
}

func (db *Bridge) GetAccountsAllowedForUser(userID string) ([]mdl.AccountMeta, error) {
	return nil, errx.ErrNotImplemented(Scope, "GetAccountsAllowedForUser")
}

func (db *Bridge) GetUserAccount(userID string) ([]mdl.Account, error) {
	return nil, errx.ErrNotImplemented(Scope, "GetUserAccount")
}

func (db *Bridge) RegisterAccount(userID string, currency string, name string) error {
	return errx.ErrNotImplemented(Scope, "RegisterAccount")
}

func (db *Bridge) GetAccount(accountID string) (*mdl.Account, error) {
	return nil, errx.ErrNotImplemented(Scope, "GetAccount")
}
