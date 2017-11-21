package model

import (
	db "github.com/gavrilaf/spawn/pkg/dbx/model"
)

type Session struct {
	ID                string
	RefreshToken      string
	ClientID          string
	ClientSecret      []byte
	UserID            string
	DeviceID          string
	IsDeviceConfirmed bool
	Locale            string
	Lang              string
	db.Permissions
}

type AuthUser struct {
	ID string
	db.AuthInfo
}

type AuthDevice struct {
	DeviceID    string
	UserID      string
	Fingerpring []byte
	IsConfirmed bool
	Locale      string
	Lang        string
}

func CreateAuthUserFromProfile(p db.UserProfile) AuthUser {
	return AuthUser{
		ID:       p.ID,
		AuthInfo: p.AuthInfo}
}

func CreateAuthDeviceFromDevice(d db.DeviceInfo) AuthDevice {
	return AuthDevice{
		DeviceID:    d.ID,
		UserID:      d.UserID,
		IsConfirmed: d.IsConfirmed,
		Fingerpring: d.Fingerprint,
		Locale:      d.Locale,
		Lang:        d.Lang}
}
