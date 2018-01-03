package mdl

import (
	"time"

	"github.com/fatih/structs"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
)

type UserDeviceInfo struct {
	ID          string `structs:"device_id"`
	Name        string `structs:"device_name"`
	IsCurrent   bool   `structs:"is_current"`
	IsConfirmed bool   `structs:"is_confirmed"`
	Locale      string `structs:"locale"`
	Lang        string `structs:"lang"`
	LoginTime   int64  `structs:"login_time"`
	LoginIP     string `structs:"login_ip"`
	UserAgent   string `structs:"user_agent"`
	LoginRegion string `structs:"login_region"`
}

func CreateUserDeviceInfoFromDb(d db.DeviceInfoEx) UserDeviceInfo {
	return UserDeviceInfo{
		ID:          d.ID,
		Name:        d.Name,
		IsCurrent:   false,
		IsConfirmed: d.IsConfirmed,
		Locale:      d.Locale,
		Lang:        d.Lang,
		LoginTime:   d.GetLoginTime().Unix(),
		LoginIP:     d.GetLoginIP(),
		UserAgent:   d.GetUserAgent(),
		LoginRegion: d.GetLoginRegion(),
	}
}

func (p UserDeviceInfo) GetLoginTime() time.Time {
	return time.Unix(p.LoginTime, 0).UTC()
}

func (p UserDeviceInfo) ToMap() map[string]interface{} {
	pm := structs.Map(p)
	pm["login_time"] = p.GetLoginTime().Format(time.RFC3339)
	return pm
}
