package auth

import (
	"time"
)

const (
	Realm            = "Crypt-Auth"
	TokenHeadName    = "Bearer"
	SigningAlgorithm = "HS256"
	TokenLookup      = "Authorization"

	SessionIDName = "session_id"
	ClientIDName  = "client_id"
	UserIDName    = "user_id"
)

type LoginParcel struct {
	ClientID       string `json:"client_id" binding:"required"`
	DeviceID       string `json:"device_id" binding:"required"`
	Username       string `json:"username" binding:"required"`
	SignedPassword string `json:"signed_password" binding:"required"`
	Signature      string `json:"signature" binding:"required"`
}

type RegisterParcel struct {
	ClientID  string `json:"client_id" binding:"required"`
	DeviceID  string `json:"device_id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type RefreshParcel struct {
	AuthToken    string `json:"auth_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenParcel struct {
	AuthToken    string
	RefreshToken string
	Expire       time.Time
}

////////////////////////////////////////////////////////////////////////

func (p *LoginParcel) CheckSignature(key []byte) bool {
	return true
}

func (p *LoginParcel) CheckPassword(password string) bool {
	return true
}

func (p *LoginParcel) CheckDevice(devices []string) bool {
	for _, d := range devices {
		if p.DeviceID == d {
			return true
		}
	}
	return false
}
