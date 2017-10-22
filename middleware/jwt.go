package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gavrilaf/go-auth/storage"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

const (
	Realm            = "Crypt-Auth"
	TokenHeadName    = "Bearer"
	SigningAlgorithm = "HS256"
	TokenLookup      = "Authorization"
)

var SecretKey = []byte("This-is-secret-key")

type Login struct {
	ClientID   string `json:"client_id" binding:"required"`
	Username   string `json:"username" binding:"required"`
	DeviceID   string `json:"device_id" binding:"required"`
	SignSecret string `json:"sign_secret" binding:"required"`
	SignKey    string `json:"sign_key" binding:"required"`
}

type TokenDesc struct {
	TokenString string
	Expire      time.Time
}

//type

// AuthMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
type AuthMiddleware struct {
	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	Storage storage.StorageFacade
}

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *AuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := mw.parseToken(c)

		if err != nil {
			mw.Unauthorized(c, http.StatusUnauthorized, err.Error())
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		sessionId := claims["session_id"].(string)
		c.Set("JWT_PAYLOAD", claims)
		c.Set("sessionId", sessionId)

		if !mw.CheckAccess(sessionId, c) {
			mw.Unauthorized(c, http.StatusForbidden, "You don't have permission to access.")
			return
		}

		c.Next()
	}
}

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *AuthMiddleware) LoginHandler(c *gin.Context) {

	var loginVals Login

	err := c.ShouldBindWith(&loginVals, binding.JSON)
	if err != nil {
		mw.Unauthorized(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := mw.HandleLogin(&loginVals)
	if err != nil {
		mw.Unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_token": token.TokenString,
		"expire":     token.Expire.Format(time.RFC3339),
	})
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the GinJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
/*func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {
	token, _ := mw.parseToken(c)
	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		mw.unauthorized(c, http.StatusUnauthorized, "Token is expired.")
		return
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	newClaims["id"] = claims["id"]
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = origIat

	tokenString, err := newToken.SignedString(mw.Key)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}*/

func (mw *AuthMiddleware) Unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm="+Realm)
	c.JSON(code, gin.H{"code": code, "message": message})
	c.Abort()
}

func (mw *AuthMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
	token, err := mw.jwtFromHeader(c, TokenLookup)
	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(jwt.MapClaims)
		aud := claims["aud"].(string)
		fmt.Printf("Inside parse token, aud = %v\n", aud)

		client, err := mw.Storage.FindClientByID(aud)
		if err != nil {
			return nil, err
		}

		return []byte(client.Secret), nil
	})
}

func (mw *AuthMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", errors.New("auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == TokenHeadName) {
		return "", errors.New("invalid auth header")
	}

	return parts[1], nil
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *gin.Context) jwt.MapClaims {

	if _, exists := c.Get("JWT_PAYLOAD"); !exists {
		emptyClaims := make(jwt.MapClaims)
		return emptyClaims
	}

	jwtClaims, _ := c.Get("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)
}
