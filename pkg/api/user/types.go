package user

import (
	"github.com/fatih/structs"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	rdm "github.com/gavrilaf/spawn/pkg/cache/model"
)

type ConfirmDeviceRequest struct {
	Code string `json:"code" form:"code" binding:"required"`
}

type DeleteDeviceRequest struct {
	DeviceID string `json:"device_id" form:"device_id" binding:"required"`
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
	Devices []rdm.UserDeviceInfo `structs:"devices"`
}

func (s UserDevices) ToMap() map[string]interface{} {
	return structs.Map(s)
}
