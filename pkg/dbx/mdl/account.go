package mdl

import "time"

type AccountMeta struct {
	Currency        string `db:"currency" structs:"currency"`
	IsCrypto        string `db:"is_crypto" structs:"is_crypto"`
	Precision       int    `db:"precision" structs:"precision"`
	MultipleAllowed bool   `db:"multiple_allowed" structs:"multiple_allowed"`
	Description     string `db:"description" structs:"description"`
}

type Account struct {
	ID       string    `db:"id" structs:"id"`
	UserID   string    `db:"user_id" structs:"user_id"`
	Name     string    `db:"name" structs:"name"`
	Currency string    `db:"currency" structs:"currency"`
	Created  time.Time `db:"created" structs:"created"`
}
