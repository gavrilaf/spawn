package auth

import (
	"fmt"
	"time"

	"github.com/gavrilaf/spawn/pkg/cryptx"
)

const (
	Realm = "Spawn"

	TokenHeadName    = "Bearer"
	SigningAlgorithm = "HS256"
	TokenLookup      = "Authorization"

	SessionIDName = "session_id"
	ClientIDName  = "client_id"
	UserIDName    = "user_id"

	AuthTypeSimple = "simple"
)

type LoginParcel struct {
	ClientID  string `json:"client_id" form:"client_id" binding:"required"`
	DeviceID  string `json:"device_id" form:"device_id" binding:"required"`
	AuthType  string `json:"auth_type" form:"auth_type" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
	Signature string `json:"signature" form:"signature" binding:"required"`
}

type RegisterParcel struct {
	ClientID  string `json:"client_id" form:"client_id" binding:"required"`
	DeviceID  string `json:"device_id" form:"device_id" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required,ascii,min=8"`
	Signature string `json:"signature" form:"signature" binding:"required"`
}

type RefreshParcel struct {
	AuthToken    string `json:"auth_token" form:"auth_token" binding:"required"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type TokenParcel struct {
	AuthToken    string
	RefreshToken string
	Expire       time.Time
}

////////////////////////////////////////////////////////////////////////

func (p *LoginParcel) CheckSignature(key []byte) error {
	msg := p.ClientID + p.DeviceID + p.Username
	return cryptx.CheckSignature(msg, p.Signature, key)
}

func (p *LoginParcel) CheckPassword(pswHash string) bool {
	return cryptx.CheckPassword(p.Password, pswHash) == nil
}

func (p *LoginParcel) CheckDevice(devices []string) bool {
	for _, d := range devices {
		if p.DeviceID == d {
			return true
		}
	}
	return false
}

func (p *LoginParcel) String() string {
	return fmt.Sprintf("LoginParcel(%v, %v, %v, %v, %v)", p.ClientID, p.DeviceID, p.AuthType, p.Username, p.Signature)
}

////////////////////////////////////////////////////////////////////////

func (p *RegisterParcel) CheckSignature(key []byte) error {
	msg := p.ClientID + p.DeviceID + p.Username
	return cryptx.CheckSignature(msg, p.Signature, key)
}

func (p *RegisterParcel) String() string {
	return fmt.Sprintf("RegisterParcel(%v, %v, %v, %v)", p.ClientID, p.DeviceID, p.Username, p.Signature)
}
