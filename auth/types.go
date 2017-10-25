package auth

import (
	"github.com/gavrilaf/go-auth/auth/storage"
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

// AuthMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
type AuthMiddleware struct {
	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	MaxRefresh time.Duration

	Storage storage.StorageFacade
}

type LoginParcel struct {
	ClientID       string `json:"client_id" binding:"required"`
	DeviceID       string `json:"device_id" binding:"required"`
	Username       string `json:"username" binding:"required"`
	SignedSecret   string `json:"signed_secret" binding:"required"`
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
