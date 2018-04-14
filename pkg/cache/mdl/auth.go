package mdl

import (
	"fmt"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
)

type Session struct {
	ID                string
	Nonce             int64
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

func (p Session) String() string {
	return fmt.Sprintf("Session{ID: %s, Nonce: %d, Client: %s, User: %s, Device: %s, DeviceConfirmed: %t, Permissions: %s, Loc: %s, Lang: %s}",
		p.ID, p.Nonce, p.ClientID, p.UserID, p.DeviceID, p.IsDeviceConfirmed, p.Permissions.String(), p.Locale, p.Lang)
}

type AuthUser struct {
	ID string
	db.AuthInfo
}

type AuthDevice struct {
	DeviceID    string
	UserID      string
	Fingerprint []byte
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
		DeviceID:    d.DeviceID,
		UserID:      d.UserID,
		IsConfirmed: d.IsConfirmed,
		Fingerprint: d.Fingerprint,
		Locale:      d.Locale,
		Lang:        d.Lang}
}
