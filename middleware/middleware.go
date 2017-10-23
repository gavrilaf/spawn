package middleware

import (
	"fmt"
	"github.com/gavrilaf/go-auth/storage"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"time"
)

func (mw *AuthMiddleware) HandleLogin(p *LoginParcel) (*TokenParcel, error) {

	client, err := mw.Storage.FindClientByID(p.ClientID)
	if err != nil {
		return nil, err
	}

	user, err := mw.Storage.FindUserByEmail(p.Username)
	if err != nil {
		return nil, err
	}

	sessionId := mw.GenerateSessionID()
	refreshToken := mw.GenerateRefreshToken(sessionId)

	session := storage.Session{ID: sessionId, RefreshToken: refreshToken, ClientID: client.ID, UserID: user.ID, Email: user.Email, Secret: client.Secret}

	err = mw.Storage.StoreSession(session)
	if err != nil {
		return nil, err
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	now := time.Now()
	expire := now.Add(mw.Timeout)
	claims["session_id"] = session.ID
	claims["aud"] = session.ClientID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = now.Unix()
	//claims["iss"] = "go-auth" // TODO: Fix it later

	tokenString, err := token.SignedString([]byte(session.Secret))
	if err != nil {
		return nil, err

	}

	return &TokenParcel{AuthToken: tokenString, RefreshToken: refreshToken, Expire: expire}, nil
}

func (mw *AuthMiddleware) HandleRefresh(p *RefreshParcel) (*TokenParcel, error) {
	token, _ := mw.parseToken(p.AuthToken)
	claims := token.Claims.(jwt.MapClaims)

	sessionId := claims["session_id"].(string)
	origIat := int64(claims["orig_iat"].(float64))

	if origIat < time.Now().Add(-mw.MaxRefresh).Unix() {
		return nil, fmt.Errorf("Token is expired.")
	}

	session, err := mw.Storage.FindSessionByID(sessionId)
	if err != nil {
		return nil, err
	}

	if p.RefreshToken != session.RefreshToken {
		return nil, fmt.Errorf("Invalid refresh token")
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	now := time.Now()
	expire := now.Add(mw.Timeout)
	claims["exp"] = expire.Unix()

	tokenString, err := token.SignedString([]byte(session.Secret))
	if err != nil {
		return nil, err
	}

	return &TokenParcel{AuthToken: tokenString, RefreshToken: "", Expire: expire}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *AuthMiddleware) CheckAccess(userId string, c *gin.Context) bool {
	return true
}

func (mw *AuthMiddleware) GenerateSessionID() string {
	return uuid.NewV4().String()
}

func (mw *AuthMiddleware) GenerateRefreshToken(sessionId string) string {
	return sessionId
}
