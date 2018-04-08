package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/api/ginx"
)

type Auth struct {
	*api.Bridge
}

func CreateAuth(bridge *api.Bridge) *Auth {
	return &Auth{Bridge: bridge}
}

func (self Auth) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := ginx.JwtFromHeader(c, defs.TokenLookup)
		if err != nil {
			ginx.HandleAuthError(c, http.StatusBadRequest, err)
		}

		token, err := ginx.ParseToken(tokenStr, func(id string) (interface{}, error) {
			cl, err := self.ReadModel.FindClient(id)
			if err != nil {
				return nil, err
			}
			return cl.Secret, nil
		})

		if err != nil {
			ginx.HandleAuthError(c, http.StatusUnauthorized, err)
			return
		}

		claims := ginx.ClaimsFromToken(token)
		session, err := self.ReadModel.GetSession(claims.SessionID())
		if err != nil {
			ginx.HandleAuthError(c, http.StatusUnauthorized, defs.ErrSessionNotFound)
			return
		}

		c.Set(defs.SessionKey, session)

		c.Next()
	}
}
