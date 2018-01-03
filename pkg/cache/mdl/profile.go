package mdl

import (
	"time"

	"github.com/fatih/structs"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
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
