package auth

import (
	//"fmt"
	"github.com/gavrilaf/spawn/pkg/errx"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/dgrijalva/jwt-go.v3"

	"net/http"
	"strings"
	"time"
)

//type
type Middleware struct {
	Timeout    time.Duration
	MaxRefresh time.Duration

	Storage StorageFacade

	Log *logrus.Logger
}

// MiddlewareFunc makes AuthMiddleware implement the Middleware interface.
func (mw *Middleware) MiddlewareFunc() gin.HandlerFunc {
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
		c.Set(DeviceIDName, session.DeviceID)

		c.Next()
	}
}

// LoginHandler can be used by clients to get a jwt token.
func (mw *Middleware) LoginHandler(c *gin.Context) {
	var loginVals LoginDTO

	err := c.Bind(&loginVals)
	if err != nil {
		mw.HandleError(c, http.StatusBadRequest, err)
		return
	}

	token, err := mw.HandleLogin(&loginVals)
	if err != nil {
		mw.HandleError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, token.ToJson())
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
func (mw *Middleware) RefreshHandler(c *gin.Context) {
	var refreshVals RefreshDTO

	err := c.Bind(&refreshVals)
	if err != nil {
		mw.HandleError(c, http.StatusUnauthorized, err)
		return
	}

	token, err := mw.HandleRefresh(&refreshVals)
	if err != nil {
		mw.HandleError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, token.ToJson())
}

func (mw *Middleware) RegisterHandler(c *gin.Context) {
	var registerVals RegisterDTO

	err := c.Bind(&registerVals)
	if err != nil {
		mw.HandleError(c, http.StatusBadRequest, err)
		return
	}

	registered, err := mw.HandleRegister(&registerVals)
	if err != nil {
		mw.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	mw.Log.Infof("User %v registered", registerVals.Username)
	c.JSON(http.StatusOK, registered.ToJson())
}

////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *Middleware) HandleError(c *gin.Context, httpCode int, err error) {
	c.Header("WWW-Authenticate", "JWT realm="+Realm)

	mw.Log.Errorf("auth.HandleError, code=%d, err=%v", httpCode, err)
	errJson := errx.Error2Json(err, errScope)

	c.JSON(httpCode, gin.H{"error": errJson})
	c.Abort()
}

////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *Middleware) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		claims := ClaimsFromToken(token)
		client, err := mw.Storage.FindClientByID(claims.ClientID())
		if err != nil {
			return nil, err
		}

		return client.Secret(), nil
	})
}

func (mw *Middleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
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
