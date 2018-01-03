package account

import (
	"github.com/fatih/structs"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
)

type UserAccounts struct {
	Accounts []mdl.Account `structs:"accounts"`
}

func (s UserAccounts) ToMap() map[string]interface{} {
	return structs.Map(s)
}
