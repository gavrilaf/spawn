package middleware

import (
	"github.com/gavrilaf/go-auth/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"net/http"
	"strings"
	"time"
)

//type

// MiddlewareFunc makes AuthMiddleware implement the Middleware interface.
func (mw *AuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := mw.jwtFromHeader(c, TokenLookup)
		if err != nil {
			mw.HandleError(c, http.StatusBadRequest, err)
		}

		token, err := mw.parseToken(tokenStr)
		if err != nil {
			mw.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		claims := ClaimsFromToken(token)
		session, err := mw.Storage.FindSessionByID(claims.SessionID())
		if err != nil {
			mw.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		if !mw.CheckAccess(session.UserID, session.ClientID, c) {
			mw.HandleError(c, http.StatusForbidden, errAccessForbiden)
			return
		}

		c.Set(ClientIDName, session.ClientID)
		c.Set(UserIDName, session.UserID)

		c.Next()
	}
}

// LoginHandler can be used by clients to get a jwt token.
func (mw *AuthMiddleware) LoginHandler(c *gin.Context) {
	var loginVals LoginParcel

	err := c.ShouldBindWith(&loginVals, binding.JSON)
	if err != nil {
		mw.HandleError(c, http.StatusBadRequest, err)
		return
	}

	token, err := mw.HandleLogin(&loginVals)
	if err != nil {
		mw.HandleError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_token":    token.AuthToken,
		"refresh_token": token.RefreshToken,
		"expire":        token.Expire.Format(time.RFC3339),
	})
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
func (mw *AuthMiddleware) RefreshHandler(c *gin.Context) {
	var refreshVals RefreshParcel

	err := c.ShouldBindWith(&refreshVals, binding.JSON)
	if err != nil {
		mw.HandleError(c, http.StatusUnauthorized, err)
		return
	}

	token, err := mw.HandleRefresh(&refreshVals)
	if err != nil {
		mw.HandleError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_token": token.AuthToken,
		"expire":     token.Expire.Format(time.RFC3339),
	})
}

func (mw *AuthMiddleware) RegisterHandler(c *gin.Context) {
	var registerVals RegisterParcel

	err := c.ShouldBindWith(&registerVals, binding.JSON)
	if err != nil {
		mw.HandleError(c, http.StatusBadRequest, err)
		return
	}

	if err = mw.HandleRegister(&registerVals); err != nil {
		mw.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *AuthMiddleware) HandleError(c *gin.Context, httpCode int, err error) {
	c.Header("WWW-Authenticate", "JWT realm="+Realm)

	var errJson map[string]interface{}

	switch err2 := err.(type) {
	case errors.Err:
		errJson = gin.H{"scope": err2.Scope, "reason": err2.Reason}
	default:
		errJson = gin.H{"scope": ErrScope, "reason": ErrReasonDefault, "message": err.Error()}
	}

	c.JSON(httpCode, gin.H{"error": errJson})
	c.Abort()
}

////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *AuthMiddleware) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		claims := ClaimsFromToken(token)
		client, err := mw.Storage.FindClientByID(claims.ClientID())
		if err != nil {
			return nil, err
		}

		return client.Secret, nil
	})
}

func (mw *AuthMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", errInvalidRequest
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == TokenHeadName) {
		return "", errInvalidRequest
	}

	return parts[1], nil
}
