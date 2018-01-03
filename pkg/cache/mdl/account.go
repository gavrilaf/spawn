package mdl

import (
	"time"

	"github.com/fatih/structs"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
)

type AccountStatus int

const (
	ACCOUNT_CREATED AccountStatus = iota
	ACCOUNT_ACTIVE
	ACCOUNT_SUSPENDED
	ACCOUNT_CLOSED
)

type AccountState struct {
	Status  AccountStatus `structs:"status"`
	Balance string        `structs:"balance"`
	Updated int64         `structs:"updated"`
}

func (p AccountState) GetUpdated() time.Time {
	return time.Unix(p.Updated, 0).UTC()
}

func (s AccountState) ToMap() map[string]interface{} {
	p := structs.Map(s)
	p["updated"] = s.GetUpdated().Format(time.RFC3339)
	return p
}

///////////////////////////////////////////////////////////////////////////////

type Account struct {
	ID       string `structs:"id"`
	Name     string `structs:"name"`
	Currency string `structs:"currency"`
	Created  int64  `structs:"created"`
	AccountState
}

func (p Account) GetCreated() time.Time {
	return time.Unix(p.Created, 0).UTC()
}

func (s Account) ToMap() map[string]interface{} {
	p := structs.Map(s)
	p["created"] = s.GetCreated().Format(time.RFC3339)
	return p
}

func CreateAcountFromDbModel(p db.Account) Account {
	return Account{
		ID:           p.ID,
		Name:         p.Name,
		Currency:     p.Currency,
		Created:      p.Created.Unix(),
		AccountState: AccountState{},
	}
}
