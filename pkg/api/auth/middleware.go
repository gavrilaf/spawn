package auth

import (
	"github.com/gavrilaf/spawn/pkg/api"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/dgrijalva/jwt-go.v3"

	"net/http"
	"strings"
	"time"
)

//type
type Middleware struct {
	timeout    time.Duration
	maxRefresh time.Duration
	bridge     Bridge
}

func CreateMiddleware(bridge *api.Bridge) *Middleware {
	if bridge == nil {
		return nil
	}

	return &Middleware{
		timeout:    time.Hour,
		maxRefresh: time.Hour * 24,
		bridge:     Bridge{bridge},
	}
}

// MiddlewareFunc makes AuthMiddleware implement the Middleware interface.
func (mw *Middleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := mw.jwtFromHeader(c, TokenLookup)
		if err != nil {
			mw.handleError(c, http.StatusBadRequest, err)
		}

		token, err := mw.parseToken(tokenStr)
		if err != nil {
			mw.handleError(c, http.StatusUnauthorized, err)
			return
		}

		claims := ClaimsFromToken(token)
		session, err := mw.bridge.GetSession(claims.SessionID())
		if err != nil {
			mw.handleError(c, http.StatusUnauthorized, err)
			return
		}

		if !mw.CheckAccess(session.UserID, session.ClientID, c) {
			mw.handleError(c, http.StatusForbidden, api.ErrAccessForbiden)
			return
		}

		c.Set("session_id", session.ID)
		c.Set("client_id", session.ClientID)
		c.Set("user_id", session.UserID)
		c.Set("device_id", session.DeviceID)

		c.Next()
	}
}

// LoginHandler can be used by clients to get a jwt token.
func (mw *Middleware) LoginHandler(c *gin.Context) {
	var loginVals LoginDTO

	err := c.Bind(&loginVals)
	if err != nil {
		mw.handleError(c, http.StatusBadRequest, err)
		return
	}

	token, err := mw.HandleLogin(loginVals, createLoginContext(c))
	if err != nil {
		mw.handleError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, token.ToMap())
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
func (mw *Middleware) RefreshHandler(c *gin.Context) {
	var refreshVals RefreshDTO

	err := c.Bind(&refreshVals)
	if err != nil {
		mw.handleError(c, http.StatusUnauthorized, err)
		return
	}

	token, err := mw.HandleRefresh(refreshVals)
	if err != nil {
		mw.handleError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, token.ToMap())
}

// RegisterHandler can be used by clients to register and get a jwt token.
func (mw *Middleware) RegisterHandler(c *gin.Context) {
	var registerVals RegisterDTO

	err := c.Bind(&registerVals)
	if err != nil {
		mw.handleError(c, http.StatusBadRequest, err)
		return
	}

	token, err := mw.HandleRegister(registerVals, createLoginContext(c))
	if err != nil {
		mw.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, token.ToMap())
}

////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *Middleware) handleError(c *gin.Context, httpCode int, err error) {
	c.Header("WWW-Authenticate", "JWT realm="+Realm)
	log.Errorf("auth error, code=%d, err=%v", httpCode, err)
	c.JSON(httpCode, gin.H{"error": errx.Error2Map(err, api.ErrScope)})
	c.Abort()
}

////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *Middleware) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		claims := ClaimsFromToken(token)
		client, err := mw.bridge.GetClient(claims.ClientID())
		if err != nil {
			return nil, err
		}

		return client.Secret, nil
	})
}

func (mw *Middleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", api.ErrInvalidRequest
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == TokenHeadName) {
		return "", api.ErrInvalidRequest
	}

	return parts[1], nil
}

////////////////////////////////////////////////////////////////////////////////////////////////
func (mw *Middleware) getClient(id string) (*db.Client, error) {
	return mw.bridge.GetClient(id)
}

func createLoginContext(c *gin.Context) LoginContext {
	return LoginContext{
		IP:        c.ClientIP(),
		UserAgent: c.Request.Header.Get("User-Agent"),
	}
}
