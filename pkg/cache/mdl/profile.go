package mdl

import (
	"github.com/fatih/structs"
	"time"

	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/utils"
)

type PersonalInfo struct {
	Country        string `structs:"country"`
	FirstName      string `structs:"first_name"`
	LastName       string `structs:"last_name"`
	BirthDate      int64  `structs:"birth_date"`
	db.PhoneNumber `structs:"phone_number"`
}

type UserProfile struct {
	ID           string `structs:"id"`
	db.AuthInfo  `structs:"auth_info"`
	PersonalInfo `structs:"personal_info"`
}

func (p UserProfile) GetBirthDate() time.Time {
	return utils.Unix2Time(p.BirthDate)
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
			m["birth_date"] = utils.FormatServerDate(p.GetBirthDate())
			pm["personal_info"] = m
		}
	}

	return pm
}
