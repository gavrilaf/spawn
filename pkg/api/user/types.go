package user

import (
	"github.com/fatih/structs"
	"github.com/gavrilaf/spawn/pkg/api/auth"
)

type UserState struct {
	UserID      string              `structs:"user_id"`
	Locale      string              `structs:"locale"`
	Lang        string              `structs:"lang"`
	Permissions auth.PermissionsDTO `structs:"permissions"`
}

func (s UserState) ToMap() map[string]interface{} {
	return structs.Map(s)
}
