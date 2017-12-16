package model

import (
	"github.com/fatih/structs"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"time"
)

type PersonalInfo struct {
	Country   string `structs:"country"`
	FirstName string `structs:"first_name"`
	LastName  string `structs:"last_name"`
	BirthDate int64  `structs:"birth_date"`
	db.PhoneNumber
}

type UserProfile struct {
	ID           string `structs:"id"`
	db.AuthInfo  `structs:"auth_info"`
	PersonalInfo `structs:"personal_info"`
}

func (p UserProfile) GetBirthDate() time.Time {
	return time.Unix(p.BirthDate, 0).UTC()
}

func CreateProfileFromDbModel(p db.UserProfile) UserProfile {
	return UserProfile{
		ID:       p.ID,
		AuthInfo: p.AuthInfo,
		PersonalInfo: PersonalInfo{
			Country:     p.Country,
			FirstName:   p.FirstName,
			LastName:    p.LastName,
			BirthDate:   p.BirthDate.Unix(),
			PhoneNumber: p.PhoneNumber,
		},
	}
}

func (p UserProfile) ToMap() map[string]interface{} {
	pm := structs.Map(p)

	auth, ok := pm["auth_info"]
	if ok {
		switch m := auth.(type) {
		case map[string]interface{}:
			delete(m, "password")
			pm["auth_info"] = m
		}
	}

	personal, ok := pm["personal_info"]
	if ok {
		switch m := personal.(type) {
		case map[string]interface{}:
			m["birth_date"] = p.GetBirthDate().Format(time.RFC3339)
			pm["personal_info"] = m
		}
	}

	return pm
}

/////////////////////////////////////////////////////////////////////////////////////////

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
