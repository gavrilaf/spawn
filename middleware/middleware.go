package middleware

import (
	//"fmt"
	"github.com/gavrilaf/go-auth/storage"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"time"
)

func (mw *AuthMiddleware) HandleLogin(p *Login) (*TokenDesc, error) {

	client, err := mw.Storage.FindClientByID(p.ClientID)
	if err != nil {
		return nil, err
	}

	user, err := mw.Storage.FindUserByEmail(p.Username)
	if err != nil {
		return nil, err
	}

	sessionId := mw.GenerateSessionID()

	// TODO: Refactor
	session := storage.Session{ID: sessionId, ClientID: client.ID, UserID: user.ID, Email: user.Email, Secret: client.Secret}

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

	return &TokenDesc{TokenString: tokenString, Expire: expire}, nil
}

func (mw *AuthMiddleware) CheckAccess(userId string, c *gin.Context) bool {
	return true
}

func (mw *AuthMiddleware) GenerateSessionID() string {
	return uuid.NewV4().String()
}
