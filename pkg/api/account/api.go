package account

import (
	"net/http"

	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/api/ginx"
	"github.com/gavrilaf/spawn/pkg/errx"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (p ApiImpl) GetAccounts(c *gin.Context) {
	session, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("AccountsApi.GetAccounts, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	userID := session.UserID

	accounts, err := p.ReadModel.GetUserAccounts(userID)
	if err != nil {
		log.Errorf("AccountsApi.GetAccounts, could not read accounts: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, err)
		return
	}

	userAccounts := UserAccounts{Accounts: accounts}

	c.JSON(http.StatusOK, userAccounts.ToMap())
}

func (p ApiImpl) GetAccountState(c *gin.Context) {
	_, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("AccountsApi.GetAccountState, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, errx.ErrNotImplemented(defs.ErrScope, "GetAccountState"))
}

func (p ApiImpl) RegisterAccount(c *gin.Context) {
	_, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("AccountsApi.RegisterAccount, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, errx.ErrNotImplemented(defs.ErrScope, "RegisterAccount"))
}

func (p ApiImpl) SuspendAccount(c *gin.Context) {
	_, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("AccountsApi.SuspendAccount, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, errx.ErrNotImplemented(defs.ErrScope, "SuspendAccount"))
}

func (p ApiImpl) ResumeAccount(c *gin.Context) {
	_, err := ginx.GetContextSession(c)
	if err != nil {
		log.Errorf("AccountsApi.ResumeAccount, could not find session: %v", err)
		ginx.HandleError(c, defs.ErrScope, http.StatusUnauthorized, err)
		return
	}

	ginx.HandleError(c, defs.ErrScope, http.StatusInternalServerError, errx.ErrNotImplemented(defs.ErrScope, "ResumeAccount"))
}
