package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gavrilaf/spawn/pkg/api/ginx"
)

// SignIn -
func (self ApiImpl) SignIn(c *gin.Context) {
	var loginVals LoginDTO

	err := c.Bind(&loginVals)
	if err != nil {
		ginx.HandleAuthError(c, http.StatusBadRequest, err)
		return
	}

	token, err := self.handleSignIn(loginVals, createLoginContext(c))
	if err != nil {
		ginx.HandleAuthError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, token.ToMap())
}

// SignUp -
func (self ApiImpl) SignUp(c *gin.Context) {
	var registerVals RegisterDTO

	err := c.Bind(&registerVals)
	if err != nil {
		ginx.HandleAuthError(c, http.StatusBadRequest, err)
		return
	}

	token, err := self.handleSignUp(registerVals, createLoginContext(c))
	if err != nil {
		ginx.HandleAuthError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, token.ToMap())
}

// RefreshToken -
func (self ApiImpl) RefreshToken(c *gin.Context) {
	var refreshVals RefreshDTO

	err := c.Bind(&refreshVals)
	if err != nil {
		ginx.HandleAuthError(c, http.StatusUnauthorized, err)
		return
	}

	token, err := self.handleRefresh(refreshVals)
	if err != nil {
		ginx.HandleAuthError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, token.ToMap())
}
