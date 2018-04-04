package user

import (
	"github.com/fatih/structs"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
)

type ConfirmDeviceCode struct {
	Code string `json:"code" form:"code" binding:"required" structs:"code"`
}

func (s ConfirmDeviceCode) ToMap() map[string]interface{} {
	return structs.Map(s)
}

////////////////////////////////////////////////////////////////////////////////////////////

type UserState struct {
	UserID      string              `structs:"user_id"`
	Locale      string              `structs:"locale"`
	Lang        string              `structs:"lang"`
	Permissions auth.PermissionsDTO `structs:"permissions"`
}

func (s UserState) ToMap() map[string]interface{} {
	return structs.Map(s)
}

////////////////////////////////////////////////////////////////////////////////////////////

type UserDevices struct {
	Devices []mdl.UserDeviceInfo `structs:"devices"`
}

func (s UserDevices) ToMap() map[string]interface{} {
	return structs.Map(s)
}
