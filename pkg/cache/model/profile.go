package model

import (
	"time"

	db "github.com/gavrilaf/spawn/pkg/dbx/model"
)

type PersonalInfo struct {
	Country   string
	FirstName string
	LastName  string
	BirthDate int64
	db.PhoneNumber
}

type UserProfile struct {
	ID string
	db.AuthInfo
	PersonalInfo
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
