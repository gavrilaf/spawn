package middleware

import (
	"github.com/gavrilaf/go-auth/storage"
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
	ClientID   string `json:"client_id" binding:"required"`
	Username   string `json:"username" binding:"required"`
	DeviceID   string `json:"device_id" binding:"required"`
	SignSecret string `json:"sign_secret" binding:"required"`
	SignKey    string `json:"sign_key" binding:"required"`
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
