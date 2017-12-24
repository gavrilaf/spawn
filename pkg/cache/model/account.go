package model

import (
	"time"

	"github.com/fatih/structs"
)

type AccountBalance struct {
	ID      string `structs:"id"`
	Balance string `structs:"balance"`
	Updated int64  `structs:"updated"`
}

func (p AccountBalance) GetUpdated() time.Time {
	return time.Unix(p.Updated, 0).UTC()
}

func (s AccountBalance) ToMap() map[string]interface{} {
	p := structs.Map(s)
	p["updated"] = s.GetUpdated().Format(time.RFC3339)
	return p
}

///////////////////////////////////////////////////////////////////////////////

type Account struct {
	ID       string `structs:"id"`
	Name     string `structs:"name"`
	Status   int    `structs:"status"`
	Currency string `structs:"currency"`
	IsCrypto bool   `structs:"is_crypto"`
	Updated  int64  `structs:"updated"`
	Created  int64  `structs:"created"`
}

func (p Account) GetUpdated() time.Time {
	return time.Unix(p.Updated, 0).UTC()
}

func (p Account) GetCreated() time.Time {
	return time.Unix(p.Updated, 0).UTC()
}

func (s Account) ToMap() map[string]interface{} {
	p := structs.Map(s)
	p["updated"] = s.GetUpdated().Format(time.RFC3339)
	p["created"] = s.GetCreated().Format(time.RFC3339)
	return p
}
