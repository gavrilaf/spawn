package user

import (
	"github.com/fatih/structs"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
)

type ConfirmDeviceRequest struct {
	Code string `json:"code" form:"code" binding:"required"`
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
