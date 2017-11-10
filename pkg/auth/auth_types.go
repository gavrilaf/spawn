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
	DeviceIDName  = "device_id"

	AuthTypeSimple = "simple"
)

type LoginDTO struct {
	ClientID  string `json:"client_id" form:"client_id" binding:"required"`
	DeviceID  string `json:"device_id" form:"device_id" binding:"required"`
	AuthType  string `json:"auth_type" form:"auth_type" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
	Signature string `json:"signature" form:"signature" binding:"required"`
}

type RegisterDTO struct {
	ClientID  string `json:"client_id" form:"client_id" binding:"required"`
	DeviceID  string `json:"device_id" form:"device_id" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required,ascii,min=8"`
	Signature string `json:"signature" form:"signature" binding:"required"`
}

type RefreshDTO struct {
	AuthToken    string `json:"auth_token" form:"auth_token" binding:"required"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type AuthTokenDTO struct {
	AuthToken    string
	RefreshToken string
	Expire       time.Time
}

type UserRegisteredDTO struct {
}

////////////////////////////////////////////////////////////////////////

func (p *LoginDTO) CheckSignature(key []byte) error {
	msg := p.ClientID + p.DeviceID + p.Username
	return cryptx.CheckSignature(msg, p.Signature, key)
}

func (p *LoginDTO) CheckPassword(pswHash string) bool {
	return cryptx.CheckPassword(p.Password, pswHash) == nil
}

func (p *LoginDTO) CheckDevice(devices []string) bool {
	for _, d := range devices {
		if p.DeviceID == d {
			return true
		}
	}
	return false
}

func (p *LoginDTO) String() string {
	return fmt.Sprintf("LoginParcel(%v, %v, %v, %v, %v)", p.ClientID, p.DeviceID, p.AuthType, p.Username, p.Signature)
}

////////////////////////////////////////////////////////////////////////

func (p *RegisterDTO) CheckSignature(key []byte) error {
	msg := p.ClientID + p.DeviceID + p.Username
	return cryptx.CheckSignature(msg, p.Signature, key)
}

func (p *RegisterDTO) String() string {
	return fmt.Sprintf("RegisterParcel(%v, %v, %v, %v)", p.ClientID, p.DeviceID, p.Username, p.Signature)
}

////////////////////////////////////////////////////////////////////////

func (p *AuthTokenDTO) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"auth_token":    p.AuthToken,
		"refresh_token": p.RefreshToken,
		"expire":        p.Expire.Format(time.RFC3339),
	}
}

////////////////////////////////////////////////////////////////////////

func (p *UserRegisteredDTO) ToJson() map[string]interface{} {
	return map[string]interface{}{"user_registered": true}
}
