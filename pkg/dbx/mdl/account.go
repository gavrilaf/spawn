package mdl

import "time"

type AccountMeta struct {
	Currency        string `db:"currency" structs:"currency"`
	IsCrypto        string `db:"is_crypto" structs:"is_crypto"`
	Precision       int    `db:"precision" structs:"precision"`
	MultipleAllowed bool   `db:"multiple_allowed" structs:"multiple_allowed"`
	LockAllowed     bool   `db:"lock_allowed" structs:"lock_allowed"`
}

type AccountStatus int

type AccountState struct {
	ID      string        `db:"id" structs:"id"`
	Balance int           `db:"balance" structs:"balance"`
	Status  AccountStatus `db:"status" structs:"status"`
	Updated time.Time     `db:"updated" structs:"updated"`
}

type Account struct {
	ID      string    `db:"id" structs:"id"`
	Name    string    `db:"name" structs:"name"`
	Created time.Time `db:"created" structs:"created"`
	AccountMeta
	AccountState
}
