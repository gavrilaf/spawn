package middleware

import (
	//"fmt"
	"github.com/gavrilaf/go-auth/storage"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func (mw *AuthMiddleware) HandleLogin(p *Login) (*storage.Session, error) {

	client, err := mw.Storage.FindClient(p.ClientID)
	if err != nil {
		return nil, err
	}

	user, err := mw.Storage.FindUserByEmail(p.Username)
	if err != nil {
		return nil, err
	}

	sessionId := mw.GenerateSessionID()

	// TODO: Refactor
	session := storage.Session{SessionID: sessionId, ClientID: client.ClientID, UserID: user.ID, Email: user.Email, Secret: client.Secret}

	err = mw.Storage.StoreSession(session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (mw *AuthMiddleware) CheckAccess(userId string, c *gin.Context) bool {
	return true
}

func (mw *AuthMiddleware) GenerateSessionID() string {
	return uuid.NewV4().String()
}
