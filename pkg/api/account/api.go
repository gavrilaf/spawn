package account

import (
	//"context"
	"github.com/gavrilaf/spawn/pkg/api"
	//"github.com/gavrilaf/spawn/pkg/api/auth"
	//pb "github.com/gavrilaf/spawn/pkg/rpc"
	"net/http"

	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (p ApiImpl) GetAccounts(c *gin.Context) {
	session, err := p.GetSession(c)
	if err != nil {
		log.Errorf("AccountsApi.GetAccounts, could not find session: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	userID := session.UserID

	accounts, err := p.ReadModel.GetUserAccounts(userID)
	if err != nil {
		log.Errorf("AccountsApi.GetAccounts, could not read accounts: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusInternalServerError, err)
		return
	}

	userAccounts := UserAccounts{Accounts: accounts}

	c.JSON(http.StatusOK, userAccounts.ToMap())
}

func (p ApiImpl) GetAccountState(c *gin.Context) {
	_, err := p.GetSession(c)
	if err != nil {
		log.Errorf("AccountsApi.GetAccountState, could not find session: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	p.HandleError(c, api.ErrScope, http.StatusInternalServerError, errx.ErrNotImplemented(api.ErrScope, "GetAccountState"))
}

func (p ApiImpl) RegisterAccount(c *gin.Context) {
	_, err := p.GetSession(c)
	if err != nil {
		log.Errorf("AccountsApi.RegisterAccount, could not find session: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	p.HandleError(c, api.ErrScope, http.StatusInternalServerError, errx.ErrNotImplemented(api.ErrScope, "RegisterAccount"))
}

func (p ApiImpl) SuspendAccount(c *gin.Context) {
	_, err := p.GetSession(c)
	if err != nil {
		log.Errorf("AccountsApi.SuspendAccount, could not find session: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	p.HandleError(c, api.ErrScope, http.StatusInternalServerError, errx.ErrNotImplemented(api.ErrScope, "SuspendAccount"))
}

func (p ApiImpl) ResumeAccount(c *gin.Context) {
	_, err := p.GetSession(c)
	if err != nil {
		log.Errorf("AccountsApi.ResumeAccount, could not find session: %v", err)
		p.HandleError(c, api.ErrScope, http.StatusUnauthorized, api.ErrSessionNotFound)
		return
	}

	p.HandleError(c, api.ErrScope, http.StatusInternalServerError, errx.ErrNotImplemented(api.ErrScope, "ResumeAccount"))
}
