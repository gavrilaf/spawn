package cache

import (
	mdl "github.com/gavrilaf/spawn/pkg/model"
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
	mdl.Permissions
}

type AuthUser struct {
	ID string
	mdl.AuthInfo
}

type AuthDevice struct {
	DeviceID    string
	UserID      string
	Fingerpring []byte
	IsConfirmed bool
	Locale      string
	Lang        string
}

func CreateAuthUserFromProfile(p mdl.UserProfile) AuthUser {
	return AuthUser{ID: p.ID, AuthInfo: p.AuthInfo}
}

func CreateAuthDeviceFromDevice(d mdl.DeviceInfo) AuthDevice {
	return AuthDevice{DeviceID: d.ID, UserID: d.UserID, IsConfirmed: d.IsConfirmed, Fingerpring: d.Fingerprint, Locale: d.Locale, Lang: d.Lang}
}
