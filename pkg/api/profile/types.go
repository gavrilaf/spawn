package profile

import (
	"github.com/fatih/structs"

	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/utils"
)

type UpdateCountryRequest struct {
	Country string `json:"country" form:"country" binding:"required"`
}

type UpdatePersonalInfoRequest struct {
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	BirthDate string `json:"birth_date" form:"birth_date" binding:"required"`
}

type AuthInfo struct {
	Username         string `structs:"username" json:"username"`
	IsLocked         bool   `structs:"is_locked" json:"is_locked"`
	IsEmailConfirmed bool   `structs:"is_email_confirmed" json:"is_email_confirmed"`
	Is2FARequired    bool   `structs:"is_2fa_required" json:"is_2fa_required"`
	Scope            int    `structs:"scope" json:"scope"`
}

type PhoneNumber struct {
	CountryCode int    `structs:"phone_country_code" structs:"phone_country_code"`
	Number      string `structs:"phone_number" structs:"phone_number"`
	IsConfirmed bool   `structs:"is_phone_confirmed" structs:"is_phone_confirmed"`
}

type PersonalInfo struct {
	Country     string `structs:"country" json:"country"`
	FirstName   string `structs:"first_name" json:"first_name"`
	LastName    string `structs:"last_name" json:"last_name"`
	BirthDate   string `structs:"birth_date" json:"birth_date"`
	PhoneNumber `structs:"phone_number" json:"phone_number"`
}

type UserProfile struct {
	ID           string `structs:"id" json:"id"`
	AuthInfo     `structs:"auth_info" json:"auth_info"`
	PersonalInfo `structs:"personal_info" json:"personal_info"`
}

func CreateAuthInfo(p db.AuthInfo) AuthInfo {
	return AuthInfo{Username: p.Username, IsLocked: p.IsLocked, IsEmailConfirmed: p.IsEmailConfirmed, Is2FARequired: p.Is2FARequired, Scope: p.Scope}
}

func CreatePhoneNumber(p db.PhoneNumber) PhoneNumber {
	return PhoneNumber{CountryCode: p.CountryCode, Number: p.Number, IsConfirmed: p.IsConfirmed}
}

func CreatePersonalInfo(p *mdl.UserProfile) PersonalInfo {
	phoneNumber := CreatePhoneNumber(p.PhoneNumber)
	birthDate := utils.FormatServerDate(p.GetBirthDate())
	return PersonalInfo{Country: p.Country, FirstName: p.FirstName, LastName: p.LastName, BirthDate: birthDate, PhoneNumber: phoneNumber}
}

func CreateUserProfile(p *mdl.UserProfile) UserProfile {
	authInfo := CreateAuthInfo(p.AuthInfo)
	personalInfo := CreatePersonalInfo(p)

	return UserProfile{ID: p.ID, AuthInfo: authInfo, PersonalInfo: personalInfo}
}

func (p UserProfile) ToMap() map[string]interface{} {
	return structs.Map(p)
}
